package controller

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"BTaskServer/util/Query"
	"BTaskServer/util/Tools"
	"BTaskServer/util/response"
	"BTaskServer/util/validatorTool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type ITransactionLogController interface {
	ApplyTransaction(c *gin.Context) // 申请提现(普通用户调用)
	DoneTransaction(c *gin.Context)  // 完成提现(管理员调用)
	GetApplyList(c *gin.Context)     // 管理员获取需要处理的提现记录
	GetMyTranList(c *gin.Context)    // 获取自己的提现记录(普通用户)
	//GetTranList(c *gin.Context)      // 获取所有提现记录(管理员)
	GetTranListByUserId(c *gin.Context) // 获取某个用户的提现记录
}

type TransactionLogController struct {
	DB *gorm.DB
}

func (t TransactionLogController) GetTranListByUserId(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	// 获取用户Id
	userid, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	// 查询该用户id是否属于该管理员
	var user1 model.User
	if res := t.DB.Where("addUserId = ? and id = ?", usermodel.ID, userid).Limit(1).Find(&user1); res.RowsAffected == 0 {
		response.Fail(c, nil, "无权查看该用户")
		return
	}

	var pageQuery Query.PageQuery
	if !validatorTool.ValidatorQuery[*Query.PageQuery](c, &pageQuery) {
		return
	}

	var transactionLogList []model.TransactionLog
	if res := t.DB.Raw("select a.* from transactionLog as a join (select id from transactionLog where userId = ? limit ?,?) b on a.id = b.id order by createTime desc", userid, ((*pageQuery.PageNum - 1) * *pageQuery.PageSize), *pageQuery.PageSize).Scan(&transactionLogList); res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	t.DB.Model(model.TransactionLog{}).Where("userId = ?", userid).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *pageQuery.PageNum,
		"result":  transactionLogList,
	}, "获取成功")
}

func (t TransactionLogController) GetMyTranList(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 0 {
		response.AuthError(c, nil, "管理员无此操作")
		return
	}

	var pageQuery Query.PageQuery
	if !validatorTool.ValidatorQuery[*Query.PageQuery](c, &pageQuery) {
		return
	}

	var transactionLogList []model.TransactionLog
	if res := t.DB.Raw("select a.* from transactionLog as a join (select id from transactionLog where userId = ? limit ?,?) b on a.id = b.id order by createTime desc", usermodel.ID, ((*pageQuery.PageNum - 1) * *pageQuery.PageSize), *pageQuery.PageSize).Scan(&transactionLogList); res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	t.DB.Model(model.TransactionLog{}).Where("userId = ?", usermodel.ID).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *pageQuery.PageNum,
		"result":  transactionLogList,
	}, "获取成功")

}

func (t TransactionLogController) GetApplyList(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	var pageQuery Query.PageQuery
	if !validatorTool.ValidatorQuery[*Query.PageQuery](c, &pageQuery) {
		return
	}

	var transactionLogList []model.TransactionLog
	if res := t.DB.Raw("select a.* from transactionLog as a join (select id from transactionLog where status = 0 and addUserId = ? limit ?,?) b on a.id = b.id order by createTime", usermodel.ID, ((*pageQuery.PageNum - 1) * *pageQuery.PageSize), *pageQuery.PageSize).Scan(&transactionLogList); res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	t.DB.Model(model.TransactionLog{}).Where("status = ? and addUserId = ?", 0, usermodel.ID).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *pageQuery.PageNum,
		"result":  transactionLogList,
	}, "获取成功")

}

func (t TransactionLogController) ApplyTransaction(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 0 {
		response.AuthError(c, nil, "管理员无此操作")
		return
	}

	// 这里要重新获取一次用户信息
	var userData model.User
	res := t.DB.Where("id = ?", usermodel.ID).Limit(1).Find(&userData)
	if res.Error != nil {
		response.ServerBad(c, nil, "查询个人信息失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, nil, "未查询到用户信息")
		return
	}

	// 检查是否设置了提现账号
	if userData.TranType == "" || userData.TranAccount == "" || userData.TranName == "" {
		response.Fail(c, nil, "未设置提现账号信息")
		return
	}

	var applyJson model.ApplyJson
	if !validatorTool.ValidatorJson[*model.ApplyJson](c, &applyJson) {
		return
	}

	minApply := viper.GetFloat64("提现配置.最低提现金额")
	// 对比最小金额的限制
	myMoney := userData.Money
	if applyJson.Applyprice < minApply {
		response.Fail(c, nil, fmt.Sprintf("最低提现金额为%f", minApply))
		return
	}

	if myMoney < applyJson.Applyprice {
		response.Fail(c, nil, "余额不足")
		return
	}

	yue := myMoney - applyJson.Applyprice

	orderId := strconv.Itoa(int(userData.ID)) + strconv.Itoa(int(Tools.GetNowUnix(true)))
	newTransactionLog := model.TransactionLog{
		UserId:      userData.ID,
		AddUserId:   userData.AddUserId,
		OrderId:     orderId,
		CreateTime:  Tools.GetDateNowFormat(true),
		Applyprice:  applyJson.Applyprice,
		OldPrice:    myMoney,
		NewPrice:    yue,
		TranType:    userData.TranType,
		TranAccount: userData.TranAccount,
		TranName:    userData.TranName,
		Status:      0,
	}

	tx := t.DB.Begin()

	if res1 := tx.Create(&newTransactionLog); res1.Error != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "提现申请失败,服务器异常")
		return
	}

	var user1 model.User
	if res2 := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userData.ID).Limit(1).Find(&user1).Update("money", yue); res2.Error != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "提现申请失败,服务器异常")
		return
	}

	tx.Commit()

	response.Success(c, nil, "提现已申请,等待管理员审核")

}

func (t TransactionLogController) DoneTransaction(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	// 获取提现Id
	TransactionId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	var doneApplyJson model.DoneApplyJson
	if !validatorTool.ValidatorJson[*model.DoneApplyJson](c, &doneApplyJson) {
		return
	}

	// 获取提现记录数据
	var transactionLog model.TransactionLog
	res := t.DB.Where("id = ? and addUserId = ?", TransactionId, usermodel.ID).Limit(1).Find(&transactionLog)
	if res.Error != nil {
		response.ServerBad(c, nil, "查询提现订单失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, nil, "你无权处理该提现申请")
		return
	}

	if transactionLog.Status != 0 {
		response.Fail(c, nil, "该提现申请已被处理")
		return
	}

	// 获取提现的用户信息
	//var user1 model.User
	//res2 := t.DB.Where("id = ?", transactionLog.UserId).Find(&user1).Limit(1)
	//if res2.Error != nil {
	//	response.ServerBad(c, nil, "查询提现用户信息失败")
	//	return
	//}
	//
	//if res2.RowsAffected == 0 {
	//	response.ServerBad(c, nil, "未查询到提现用户信息")
	//	return
	//}

	tx := t.DB.Begin()

	// 如果审核失败了,要同时把用户的金额改回去
	if doneApplyJson.Status != 1 {
		var user2 model.User
		res2 := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", transactionLog.UserId).Limit(1).Find(&user2).Update("money", gorm.Expr(fmt.Sprintf("money+%f", transactionLog.Applyprice)))
		if res2.Error != nil {
			tx.Rollback()
			response.ServerBad(c, nil, "修改数据失败")
			return
		}

		if res2.RowsAffected == 0 {
			tx.Rollback()
			response.ServerBad(c, nil, "修改数据失败")
			return
		}
	}

	res1 := tx.Model(&model.TransactionLog{}).Where("id = ?", TransactionId).Updates(map[string]interface{}{
		"status": doneApplyJson.Status,
		"desc":   doneApplyJson.Desc,
	})

	if res1.Error != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "修改数据失败")
		return
	}

	if res1.RowsAffected == 0 {
		tx.Rollback()
		response.ServerBad(c, nil, "修改数据失败")
		return
	}

	tx.Commit()

	response.Success(c, nil, "操作成功")
}

func NewTransactionLogController() ITransactionLogController {
	db := common.GetDB()
	if err := db.AutoMigrate(&model.TransactionLog{}); err != nil {
		panic("transactionLog表迁移失败")
	}

	fmt.Println("transactionLog表迁移成功")

	return TransactionLogController{DB: db}
}
