package controller

//
//import (
//	"BTaskServer/common"
//	"BTaskServer/model"
//	"BTaskServer/util/Query"
//	"BTaskServer/util/Tools"
//	"BTaskServer/util/response"
//	"BTaskServer/util/validatorTool"
//	"context"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/redis/go-redis/v9"
//	"github.com/spf13/viper"
//	"gorm.io/gorm"
//	"sync/atomic"
//	"time"
//)
//
//type ITaskListController interface {
//	//GetPlayTask(c *gin.Context)    // 领取播放任务
//	//SubmitPlayTask(c *gin.Context) // 提交播放任务
//	GetTaskList(c *gin.Context)      // 获取任务列表（管理员）
//	GetTaskListByKey(c *gin.Context) // 搜索列表(管理员)
//	GetTask(c *gin.Context)          // 领取任务
//	SubmitTask(c *gin.Context)       // 提交任务
//}
//
//var currentRequests int32
//var typeList []int = []int{11, 12}
//
//type TaskListController struct {
//	DB                    *gorm.DB
//	RDS                   *redis.Client
//	TaskTimeOut           int
//	BfMaxCount            int
//	MaxConcurrentRequests int32
//	MaxRedisTime          int
//}
//
//var ctx = context.Background()
//
//// 在查询前先通过Redis快速过滤
//func fastCheck(redis *redis.Client, account string, timeInt int) bool {
//	key := fmt.Sprintf("task:account:%s", account)
//	// 检查10秒内是否请求过
//	exists, _ := redis.Exists(ctx, key).Result()
//	if exists > 0 {
//		return false
//	}
//	// 设置临时标记（5秒过期）
//	redis.SetEx(ctx, key, 1, time.Duration(timeInt)*time.Second)
//	return true
//}
//
//func (t TaskListController) GetTask(c *gin.Context) {
//	// 检查当前请求数是否超出最大并发数
//	if atomic.LoadInt32(&currentRequests) >= t.MaxConcurrentRequests {
//		response.ServerBad(c, nil, "暂无任务")
//		return
//	}
//	atomic.AddInt32(&currentRequests, 1)        // 增加计数
//	defer atomic.AddInt32(&currentRequests, -1) // 使用 defer 确保请求结束时减少计数
//
//	var getTaskQuery model.GetTaskQuery
//	if !validatorTool.ValidatorQuery[*model.GetTaskQuery](c, &getTaskQuery) {
//		return
//	}
//
//	if fastCheck(t.RDS, getTaskQuery.BlAccount, t.MaxRedisTime) == false {
//		response.Fail(c, nil, "请求频繁")
//		return
//	}
//
//	// 验证type是否正确
//	isIn := false
//	for _, value := range typeList {
//		if value == getTaskQuery.Type {
//			isIn = true
//			break
//		}
//	}
//
//	if isIn == false {
//		response.Fail(c, nil, "type参数错误,没有此类型的任务")
//		return
//	}
//
//	// 验证ukey是否存在
//	user, b := getCacelUser(getTaskQuery.UserKey)
//	if b == false {
//		response.ServerBad(c, nil, "ukey验证失败")
//		return
//	}
//	// ukey不能是管理员的
//	if user.Authority == 1 {
//		response.AuthError(c, nil, "不能用管理员的ukey领取任务")
//		return
//	}
//
//	// 通过获取的订单类型，查询设置的单价
//	priceMap := GetPriceByType(user.AddUserId, getTaskQuery.Type)
//	if priceMap["ret"].(bool) == false {
//		response.ServerBad(c, nil, "领取失败,单价获取失败")
//		return
//	}
//
//	price := priceMap["price"].(float64)
//
//	var taskList model.TaskList
//	var whereStr string
//	whereStr = fmt.Sprintf("UStatus = %d and collectNum < buyNumber and getTaskType = %d and orderId not in (select orderId from taskLog where Account = '%s')", 0, getTaskQuery.Type, getTaskQuery.BlAccount)
//	res := t.DB.Where(whereStr).Limit(1).Find(&taskList)
//
//	//sqlStr := fmt.Sprintf("select * from taskList task where UStatus = %d and collectNum < buyNumber and getTaskType = %d and NOT EXISTS(select 1 from taskLog log where task.orderId=log.orderId and log.Account='%S') LIMIT 1",0, getTaskQuery.Type, getTaskQuery.BlAccount)
//	//res := t.DB.Raw("select * from taskList task where UStatus = ? and collectNum < buyNumber and getTaskType = ? and NOT EXISTS(select 1 from taskLog log where task.orderId=log.orderId and log.Account=?) LIMIT 1", 0, getTaskQuery.Type, getTaskQuery.BlAccount).Scan(&taskList)
//	//res := t.DB.Raw("SELECT t.* FROM taskList t LEFT JOIN taskLog log ON log.orderId = t.orderId AND log.Account = ? WHERE t.UStatus = ? AND t.collectNum < t.buyNumber AND t.getTaskType = ? AND log.orderId IS NULL LIMIT 1;", getTaskQuery.BlAccount, 0, getTaskQuery.Type).Scan(&taskList)
//	if res.Error != nil {
//		response.ServerBad(c, nil, "领取失败1")
//		return
//	}
//
//	if res.RowsAffected == 0 {
//		response.Fail(c, nil, "暂无任务")
//		return
//	}
//
//	taskList.CollectNum = taskList.CollectNum + 1
//	taskList.LastGetTime = Tools.GetDateNowFormat(true)
//	res1 := t.DB.Save(&taskList)
//	if res1.Error != nil {
//		response.ServerBad(c, nil, "领取失败3")
//		return
//	}
//
//	newTaskLog := model.TaskLog{
//		CollectTime: Tools.GetDateNowFormat(true),
//		OrderNum:    Tools.GetUuid(),
//		Account:     getTaskQuery.BlAccount,
//		VideoLink:   taskList.VideoLink,
//		OrderId:     taskList.OrderId,
//		GoodsId:     taskList.GoodsId,
//		GoodsName:   taskList.GoodsName,
//		Price:       price,
//		UserKey:     user.UserKey,
//		ManagerId:   user.AddUserId,
//		Status:      0,
//		GetTaskType: taskList.GetTaskType,
//		PingtaiName: taskList.PingtaiName,
//	}
//
//	if res2 := t.DB.Create(&newTaskLog); res2.Error != nil {
//		response.ServerBad(c, nil, "领取失败")
//		return
//	}
//
//	resultTask := model.GetTaskModel{
//		ID:        int(newTaskLog.ID),
//		VideoLink: taskList.VideoLink,
//		ShortLink: taskList.ShortLink,
//		VideoId:   taskList.VideoId,
//		Uid:       taskList.Uid,
//		GoodsName: taskList.GoodsName,
//	}
//
//	fmt.Println(fmt.Sprintf("%s,orderid:%d,ks号:%s,%s", taskList.PingtaiName, taskList.OrderId, getTaskQuery.BlAccount, Tools.GetDateNowFormat(true)))
//
//	response.Success(c, gin.H{"task": resultTask}, fmt.Sprintf("任务领取成功,请在%d分钟内提交", t.TaskTimeOut))
//}
//
//func (t TaskListController) GetTaskListByKey(c *gin.Context) {
//	user, _ := c.Get("user")
//	usermodel := user.(model.User)
//
//	if usermodel.Authority != 1 {
//		response.AuthError(c, nil, "权限不足")
//		return
//	}
//
//	var findTaskQuery model.FindTaskQuery
//	if !validatorTool.ValidatorQuery[*model.FindTaskQuery](c, &findTaskQuery) {
//		return
//	}
//
//	orderStr := "%" + findTaskQuery.OrderId + "%"
//
//	var taskList []model.TaskList
//	if res := t.DB.Raw("select a.* from taskList as a join (select id from taskList where orderId like ? order by orderId limit ?,?) b on a.id = b.id", orderStr, ((*findTaskQuery.PageNum - 1) * *findTaskQuery.PageSize), *findTaskQuery.PageSize).Scan(&taskList); res.Error != nil {
//		response.ServerBad(c, nil, "获取失败")
//		return
//	}
//
//	var total int64
//	t.DB.Model(model.TaskList{}).Where("orderId like ?", orderStr).Count(&total)
//
//	response.Success(c, gin.H{
//		"total":   total,
//		"pageNum": *findTaskQuery.PageNum,
//		"result":  taskList,
//	}, "获取成功")
//}
//
//func (t TaskListController) GetTaskList(c *gin.Context) {
//	user, _ := c.Get("user")
//	usermodel := user.(model.User)
//
//	if usermodel.Authority != 1 {
//		response.AuthError(c, nil, "权限不足")
//		return
//	}
//
//	var pageQuery Query.PageQuery
//	if !validatorTool.ValidatorQuery[*Query.PageQuery](c, &pageQuery) {
//		return
//	}
//	var taskList []model.TaskList
//	if res := t.DB.Raw("select a.* from taskList as a join (select id from taskList order by downTime desc limit ?,?) b on a.id = b.id", ((*pageQuery.PageNum - 1) * *pageQuery.PageSize), *pageQuery.PageSize).Scan(&taskList); res.Error != nil {
//		response.ServerBad(c, nil, "获取失败")
//		return
//	}
//
//	var total int64
//	t.DB.Model(model.TaskList{}).Where("1=1").Count(&total)
//
//	response.Success(c, gin.H{
//		"total":   total,
//		"pageNum": *pageQuery.PageNum,
//		"result":  taskList,
//	}, "获取成功")
//}
//
//func (t TaskListController) SubmitTask(c *gin.Context) {
//	var submitTaskQuery model.SubmitTaskQuery
//	if !validatorTool.ValidatorQuery[*model.SubmitTaskQuery](c, &submitTaskQuery) {
//		return
//	}
//
//	var taskLog model.TaskLog
//	res := t.DB.Where("id = ?", submitTaskQuery.ID).Limit(1).Find(&taskLog)
//	if res.Error != nil {
//		response.ServerBad(c, nil, "提交失败,服务器错误")
//		return
//	}
//
//	if res.RowsAffected == 0 {
//		response.Fail(c, nil, "提交失败,未找到领取记录")
//		return
//	}
//
//	if taskLog.Status == 2 {
//		response.Fail(c, nil, "提交失败,任务已超时")
//		return
//	}
//
//	if taskLog.Status != 0 {
//		response.Fail(c, nil, "提交失败,不能重复提交")
//		return
//	}
//
//	// 更新taskLog表
//	res1 := t.DB.Model(&model.TaskLog{}).Where("id = ?", taskLog.ID).Updates(map[string]interface{}{
//		"sumbmitTime": Tools.GetDateNowFormat(true),
//		"status":      1,
//	})
//	if res1.Error != nil {
//		response.ServerBad(c, nil, "提交失败,服务器错误")
//		return
//	}
//	if res1.RowsAffected == 0 {
//		response.ServerBad(c, nil, "提交失败,未找到领取记录")
//		return
//	}
//
//	response.Success(c, nil, "提交成功!")
//}
//
//func NewTaskListController() ITaskListController {
//	db := common.GetDB()
//
//	if err := db.AutoMigrate(&model.TaskList{}); err != nil {
//		fmt.Println(err.Error())
//		panic("taskList表迁移失败")
//	}
//	fmt.Println("taskList表迁移成功")
//
//	rds := common.GetRedis()
//
//	timeout := viper.GetInt("TaskConfig.TaskTimeOut")
//	bfMaxCount := viper.GetInt("TaskConfig.BFMaxCount")
//	maxConcurrentRequests := viper.GetInt32("system.maxConcurrentRequests")
//	maxReisTime := viper.GetInt("system.maxRedisTime")
//
//	return TaskListController{DB: db, RDS: rds, TaskTimeOut: timeout, BfMaxCount: bfMaxCount, MaxConcurrentRequests: maxConcurrentRequests, MaxRedisTime: maxReisTime}
//}
