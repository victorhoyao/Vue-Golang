package controller

import (
	"BTaskServer/common"
	"BTaskServer/global"
	"BTaskServer/model"
	"BTaskServer/util/Query"
	"BTaskServer/util/Tools"
	"BTaskServer/util/response"
	"BTaskServer/util/validatorTool"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IUserController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	MyInfo(c *gin.Context)
	AddUser(c *gin.Context)
	AddManager(c *gin.Context)
	EditUserPass(c *gin.Context)
	GetUserList(c *gin.Context)
	DelUser(c *gin.Context)
	SetTran(c *gin.Context)                // 修改提现账号(普通用户)
	ManagerSetTranByUserId(c *gin.Context) // 修改提现账号(管理员修改用户提现账号)
	FindUser(c *gin.Context)               // 模糊搜索用户
	ChangePassword(c *gin.Context)         // 修改用户密码
}

type UserController struct {
	DB *gorm.DB
}

func (u UserController) FindUser(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "无权限")
		return
	}

	var findUserQuery model.FindUserQuery
	if !validatorTool.ValidatorQuery[*model.FindUserQuery](c, &findUserQuery) {
		return
	}

	var userList []model.User
	whereStr := "%" + findUserQuery.KeyWord + "%"
	res := u.DB.Where("userName like ?", whereStr).Find(&userList)
	if res.Error != nil {
		response.ServerBad(c, nil, "服务器错误")
		return
	}

	response.Success(c, gin.H{"userList": userList}, "查询成功")
}

var UserList []model.User

func (u UserController) ManagerSetTranByUserId(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "暂无权限")
		return
	}

	userid, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	var setTranJson model.SetTranJson
	if !validatorTool.ValidatorJson(c, &setTranJson) {
		return
	}

	// 查询该管理员是否有权修改
	var user1 model.User
	if res1 := u.DB.Where("id = ? and addUserId = ?", userid, usermodel.ID).Limit(1).Find(&user1); res1.RowsAffected == 0 {
		response.Fail(c, nil, "你无权修改此用户")
		return
	}

	if res := u.DB.Model(&model.User{}).Where("id = ? and addUserId = ?", userid, usermodel.ID).Updates(map[string]interface{}{
		"tranType":    setTranJson.AccountType,
		"tranAccount": setTranJson.Account,
		"tranName":    setTranJson.Name,
	}); res.Error != nil {
		response.ServerBad(c, nil, "修改失败")
		return
	}

	cacheUserList()
	response.Success(c, nil, "修改成功")
}

func (u UserController) SetTran(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 0 {
		response.AuthError(c, nil, "管理员无此操作")
		return
	}

	// 检测是否已经修改过了
	if usermodel.TranType != "" || usermodel.TranAccount != "" || usermodel.TranName != "" {
		response.Fail(c, nil, "你已经更改过提现账号，继续修改请联系管理员")
		return
	}

	var setTranJson model.SetTranJson
	if !validatorTool.ValidatorJson(c, &setTranJson) {
		return
	}

	if res := u.DB.Model(&model.User{}).Where("id = ?", usermodel.ID).Updates(map[string]interface{}{
		"tranType":    setTranJson.AccountType,
		"tranAccount": setTranJson.Account,
		"tranName":    setTranJson.Name,
	}); res.Error != nil {
		response.ServerBad(c, nil, "设置失败")
		return
	}

	cacheUserList()
	response.Success(c, nil, "设置成功")
}

func (u UserController) DelUser(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	var delUserQuery model.DelUser
	if !validatorTool.ValidatorQuery[*model.DelUser](c, &delUserQuery) {
		return
	}

	var delUser model.User
	res := u.DB.Where("id = ? AND addUserId = ?", delUserQuery.UserId, usermodel.ID).Limit(1).Find(&delUser)
	if res.Error != nil {
		response.ServerBad(c, nil, "查询用户失败")
		return
	}
	if res.RowsAffected == 0 {
		response.AuthError(c, nil, "无权删除此用户")
		return
	}

	if res1 := u.DB.Where("id = ? AND addUserId = ?", delUserQuery.UserId, usermodel.ID).Delete(&model.User{}); res1.Error != nil {
		response.ServerBad(c, nil, "删除失败")
		return
	}

	cacheUserList()

	response.Success(c, nil, "删除成功")
}

func (u UserController) MyInfo(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	var userModel model.User
	res := u.DB.Where("id = ?", usermodel.ID).Limit(1).Find(&userModel)
	if res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	if res.RowsAffected == 0 {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	response.Success(c, gin.H{"myInfo": userModel}, "获取成功！")
}

func (u UserController) EditUserPass(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	var editUserPass model.EditUserPass
	if !validatorTool.ValidatorJson[*model.EditUserPass](c, &editUserPass) {
		return
	}

	var editUser model.User
	if res := u.DB.Where("id = ? AND addUserId = ?", editUserPass.UserId, usermodel.ID).First(&editUser); res.RowsAffected == 0 {
		response.AuthError(c, nil, "无权修改此用户")
		return
	}

	if res := u.DB.Model(&model.User{}).Where("id = ? AND addUserId = ?", editUserPass.UserId, usermodel.ID).Updates(map[string]interface{}{
		"passWord": Tools.GenMd5(editUserPass.NewPassword),
	}); res.Error != nil {
		response.ServerBad(c, nil, "修改失败")
		return
	}

	cacheUserList()

	response.Success(c, nil, "修改成功")
}

func (u UserController) AddManager(c *gin.Context) {
	var addManager model.AddManager
	if !validatorTool.ValidatorJson(c, &addManager) {
		return
	}

	newUser := model.User{
		UserName:   addManager.UserName,
		PassWord:   Tools.GenMd5(addManager.PassWord),
		Authority:  1,
		AddUserId:  0,
		UserKey:    "USER_" + strings.ToUpper(Tools.GenMd5(Tools.GetUuid())),
		Money:      float64(0.00),
		CreateTime: Tools.GetDateNowFormat(true),
	}
	if res := u.DB.Create(&newUser); res.RowsAffected == 0 {
		response.ServerBad(c, nil, "添加管理员失败")
		return
	}

	// 顺便添加Manage
	newManager := model.Manager{
		ManagerId:  newUser.ID,
		UpdateTime: Tools.GetDateNowFormat(true),
	}
	u.DB.Create(&newManager)

	cacheUserList()
	response.Success(c, nil, "添加管理员成功")

}

func (u UserController) GetUserList(c *gin.Context) {
	var pageQuery Query.PageQuery
	if !validatorTool.ValidatorQuery[*Query.PageQuery](c, &pageQuery) {
		return
	}

	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	var userlist []model.User
	if res := u.DB.Raw("select a.* from user as a join (select id from user where addUserId = ? limit ?,?) b on a.id = b.id", usermodel.ID, ((*pageQuery.PageNum - 1) * *pageQuery.PageSize), *pageQuery.PageSize).Scan(&userlist); res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	u.DB.Model(model.User{}).Where("addUserId = ?", usermodel.ID).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *pageQuery.PageNum,
		"result":  userlist,
	}, "获取成功！")

}

func (u UserController) Login(c *gin.Context) {
	// 用户登录
	var UserLogin model.UserLogin

	var UserModel model.User
	if !validatorTool.ValidatorJson[*model.UserLogin](c, &UserLogin) {
		return
	}

	// 校验用户名
	res := u.DB.Where("userName = ? AND passWord = ?", UserLogin.UserName, Tools.GenMd5(UserLogin.PassWord)).Limit(1).Find(&UserModel)
	if res.Error != nil {
		response.ServerBad(c, nil, "数据库查询失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, nil, "用户名密码错误！")
		return
	}

	// 发放token
	token, err := common.ReleseToken(UserModel)
	if err != nil {
		response.ServerBad(c, nil, "系统错误，登录失败")
		return
	}

	response.Success(c, gin.H{"token": token, "id": UserModel.ID, "authority": UserModel.Authority, "ukey": UserModel.UserKey}, "登录成功！")
}

func (u UserController) AddUser(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	var addUser model.AddUser

	if !validatorTool.ValidatorJson[*model.AddUser](c, &addUser) {
		return
	}

	// 过滤用户名重复
	var userModel model.User
	res := u.DB.Where("userName = ?", addUser.UserName).Limit(1).Find(&userModel)
	if res.Error != nil {
		response.ServerBad(c, nil, "数据库查询失败")
		return
	}
	if res.RowsAffected != 0 {
		response.Fail(c, nil, "该用户名已被占用")
		return
	}

	newUser := model.User{
		UserName:   addUser.UserName,
		PassWord:   Tools.GenMd5(addUser.PassWord),
		Authority:  0,
		AddUserId:  usermodel.ID,
		UserKey:    "USER_" + strings.ToUpper(Tools.GenMd5(Tools.GetUuid())),
		Money:      float64(0.00),
		CreateTime: Tools.GetDateNowFormat(true),
	}
	if res := u.DB.Create(&newUser); res.RowsAffected == 0 {
		response.ServerBad(c, nil, "添加用户失败")
		return
	}

	cacheUserList()

	response.Success(c, nil, "添加用户成功")

}

func cacheUserList() bool {
	db := global.GVA_DB
	UserList = []model.User{}
	res := db.Where("1=1").Find(&UserList)
	if res.Error != nil {
		fmt.Println("缓存userList失败")
		return false
	}
	return true
}

func getCacelUser(ukey string) (model.User, bool) {
	var resUser model.User
	flag := false
	for _, user := range UserList {
		if user.UserKey == ukey {
			resUser = user
			flag = true
			break
		}
	}

	return resUser, flag
}

func GetCacelUserById(id uint) (model.User, bool) {
	var resUser model.User
	flag := false
	for _, user := range UserList {
		if user.ID == id {
			resUser = user
			flag = true
			break
		}
	}

	return resUser, flag
}

func getUserList() ([]model.User, error) {
	db := global.GVA_DB
	var userList []model.User
	res := db.Where("1=1").Find(&userList)
	if res.Error != nil {
		return nil, res.Error
	}
	return userList, nil
}

func NewUserController() IUserController {
	userDB := global.GVA_DB
	if err := userDB.AutoMigrate(&model.User{}); err != nil {
		panic("user表迁移失败")
	}
	fmt.Println("user表迁移成功")

	cacheRet := cacheUserList()
	if cacheRet == false {
		panic("缓存userList失败")
	}

	return UserController{DB: userDB}
}

// Register handles public user registration
func (u UserController) Register(c *gin.Context) {
	var registerData model.AddUser
	if !validatorTool.ValidatorJson[*model.AddUser](c, &registerData) {
		return
	}

	// Check if username is already taken (only check for manager username to avoid conflict with manager's own update)
	var existingUser model.User
	if res := u.DB.Where("userName = ? AND authority = ?", registerData.UserName, 1).First(&existingUser); res.RowsAffected != 0 {
		response.Fail(c, nil, "该管理员用户名已被占用")
		return
	}

	// Start a transaction
	tx := u.DB.Begin()

	var managerUser model.User
	// Try to find an existing manager user
	if res := tx.Where("authority = ?", 1).First(&managerUser); res.Error != nil && res.RowsAffected == 0 {
		// No manager found, create the first one
		newManager := model.User{
			UserName:   registerData.UserName,
			PassWord:   Tools.GenMd5(registerData.PassWord),
			Authority:  1, // Manager
			AddUserId:  0, // No parent for the first manager
			UserKey:    "USER_" + strings.ToUpper(Tools.GenMd5(Tools.GetUuid())),
			Money:      float64(0.00),
			CreateTime: Tools.GetDateNowFormat(true),
			Status:     1,
		}
		if res := tx.Create(&newManager); res.RowsAffected == 0 {
			tx.Rollback()
			response.ServerBad(c, nil, "创建管理员失败")
			return
		}
		managerUser = newManager // Set the newly created manager as the target

	} else if res.Error != nil {
		// Other database error when trying to find manager
		tx.Rollback()
		response.ServerBad(c, nil, "查询管理员失败")
		return
	} else {
		// Manager found, update its credentials
		if res := tx.Model(&model.User{}).Where("id = ?", managerUser.ID).Updates(map[string]interface{}{
			"userName": registerData.UserName,
			"passWord": Tools.GenMd5(registerData.PassWord),
		}); res.Error != nil {
			tx.Rollback()
			response.ServerBad(c, nil, "更新管理员信息失败")
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "注册失败")
		return
	}

	// Update user cache
	cacheUserList()

	response.Success(c, nil, "管理员注册/更新成功")
}

// 修改用户密码
func (u UserController) ChangePassword(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	var changePassJson model.ChangePassJson
	if !validatorTool.ValidatorJson[*model.ChangePassJson](c, &changePassJson) {
		return
	}

	// 校验旧密码是否正确
	if usermodel.PassWord != Tools.GenMd5(changePassJson.OldPass) {
		response.Fail(c, nil, "旧密码不正确")
		return
	}

	// 更新密码
	if res := u.DB.Model(&model.User{}).Where("id = ?", usermodel.ID).Update("passWord", Tools.GenMd5(changePassJson.NewPass)); res.Error != nil {
		response.ServerBad(c, nil, "密码修改失败")
		return
	}

	// 刷新缓存
	cacheUserList()

	response.Success(c, nil, "密码修改成功")
}
