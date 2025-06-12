package controller

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"BTaskServer/util/Query"
	"BTaskServer/util/response"
	"BTaskServer/util/validatorTool"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type ITaskLogController interface {
	GetTaskLogList(c *gin.Context)      // 获取所有领取记录(管理员)
	GetTaskLogListByKey(c *gin.Context) // 搜索领取记录(管理员)
	GetMyTaskLogList(c *gin.Context)    // 获取自己的领取记录(用户)
	GetTaskLogListById(c *gin.Context)  // 获取指定用户id的领取记录(管理)
	GetTaskLogCount(c *gin.Context)     // 获取统计
}

type TaskLogController struct {
	DB *gorm.DB
}

func (t TaskLogController) GetTaskLogListByKey(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	var findTaskLogQuery model.FindTaskLogQuery
	if !validatorTool.ValidatorQuery[*model.FindTaskLogQuery](c, &findTaskLogQuery) {
		return
	}

	preDate := findTaskLogQuery.PreDate + " 00:00:00"
	nextDate := findTaskLogQuery.NextDate + " 23:59:59"

	var taskLogList []model.TaskLog
	if res := t.DB.Raw("select a.* from taskLog as a join (select id from taskLog where Account = ? and collectTime >= ? and collectTime <= ? order by collectTime DESC limit ?,?) b on a.id = b.id", findTaskLogQuery.KsAccount, preDate, nextDate, ((*findTaskLogQuery.PageNum - 1) * *findTaskLogQuery.PageSize), *findTaskLogQuery.PageSize).Scan(&taskLogList); res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	t.DB.Model(model.TaskLog{}).Where("Account = ? and collectTime >= ? and collectTime <= ?", findTaskLogQuery.KsAccount, preDate, nextDate).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *findTaskLogQuery.PageNum,
		"result":  taskLogList,
	}, "获取成功")
}

func (t TaskLogController) GetTaskLogCount(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	var getTaskLogCountQuery model.GetTaskLogCountQuery
	if !validatorTool.ValidatorQuery[*model.GetTaskLogCountQuery](c, &getTaskLogCountQuery) {
		return
	}

	date := getTaskLogCountQuery.Date + "%"

	var FailCount int64
	var SucessCount int64
	var WaitCount int64
	if usermodel.Authority == 1 {
		// 失败的
		t.DB.Model(&model.TaskLog{}).Where("collectTime like ? and getTaskType = ? and status = ? and managerId = ?", date, getTaskLogCountQuery.Type, 4, usermodel.ID).Count(&FailCount)
		// 成功的
		t.DB.Model(&model.TaskLog{}).Where("collectTime like ? and getTaskType = ? and status = ? and managerId = ?", date, getTaskLogCountQuery.Type, 3, usermodel.ID).Count(&SucessCount)
		// 等待的
		t.DB.Model(&model.TaskLog{}).Where("collectTime like ? and getTaskType = ? and status = ? and managerId = ?", date, getTaskLogCountQuery.Type, 1, usermodel.ID).Count(&WaitCount)
	} else {
		// 失败的
		t.DB.Model(&model.TaskLog{}).Where("collectTime like ? and getTaskType = ? and status = ? and userKey = ?", date, getTaskLogCountQuery.Type, 4, usermodel.UserKey).Count(&FailCount)
		// 成功的
		t.DB.Model(&model.TaskLog{}).Where("collectTime like ? and getTaskType = ? and status = ? and userKey = ?", date, getTaskLogCountQuery.Type, 3, usermodel.UserKey).Count(&SucessCount)
		// 等待的
		t.DB.Model(&model.TaskLog{}).Where("collectTime like ? and getTaskType = ? and status = ? and userKey = ?", date, getTaskLogCountQuery.Type, 1, usermodel.UserKey).Count(&WaitCount)
	}

	response.Success(c, gin.H{
		"waitCount":   WaitCount,
		"sucessCount": SucessCount,
		"failCount":   FailCount,
	}, "获取成功")

}

func (t TaskLogController) GetTaskLogListById(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	userid, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	var puser model.User
	res1 := t.DB.Where("id = ?", userid).Limit(1).Find(&puser)
	if res1.Error != nil {
		response.ServerBad(c, nil, "服务器错误")
		return
	}

	if res1.RowsAffected == 0 {
		response.Fail(c, nil, "id错误,用户不存在")
		return
	}

	var pageQuery Query.PageQuery
	if !validatorTool.ValidatorQuery[*Query.PageQuery](c, &pageQuery) {
		return
	}

	var taskLogList []model.TaskLog
	if res := t.DB.Raw("select a.* from taskLog as a join (select id from taskLog where userKey = ? order by collectTime desc limit ?,?) b on a.id = b.id", puser.UserKey, ((*pageQuery.PageNum - 1) * *pageQuery.PageSize), *pageQuery.PageSize).Scan(&taskLogList); res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	t.DB.Model(model.TaskLog{}).Where("userKey = ?", puser.UserKey).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *pageQuery.PageNum,
		"result":  taskLogList,
	}, "获取成功")

}

func (t TaskLogController) GetMyTaskLogList(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority == 1 {
		response.AuthError(c, nil, "管理员角色无结果")
		return
	}

	var pageQuery Query.PageQuery
	if !validatorTool.ValidatorQuery[*Query.PageQuery](c, &pageQuery) {
		return
	}

	var taskLogList []model.TaskLog
	if res := t.DB.Raw("select a.* from taskLog as a join (select id from taskLog where userKey = ? order by collectTime desc limit ?,?) b on a.id = b.id", usermodel.UserKey, ((*pageQuery.PageNum - 1) * *pageQuery.PageSize), *pageQuery.PageSize).Scan(&taskLogList); res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	t.DB.Model(model.TaskLog{}).Where("userKey = ?", usermodel.UserKey).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *pageQuery.PageNum,
		"result":  taskLogList,
	}, "获取成功")

}

func (t TaskLogController) GetTaskLogList(c *gin.Context) {
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

	var taskLogList []model.TaskLog
	if res := t.DB.Raw("select a.* from taskLog as a join (select id from taskLog where managerId = ? order by collectTime desc limit ?,?) b on a.id = b.id", usermodel.ID, ((*pageQuery.PageNum - 1) * *pageQuery.PageSize), *pageQuery.PageSize).Scan(&taskLogList); res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	t.DB.Model(model.TaskLog{}).Where("managerId = ?", usermodel.ID).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *pageQuery.PageNum,
		"result":  taskLogList,
	}, "获取成功")
}

func NewTaskLogController() ITaskLogController {
	db := common.GetDB()
	if err := db.AutoMigrate(&model.TaskLog{}); err != nil {
		panic("taskLog表迁移失败")
	}
	fmt.Println("taskLog表迁移成功")
	return TaskLogController{DB: db}
}
