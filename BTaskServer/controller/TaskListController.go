package controller

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"BTaskServer/util/Query"
	"BTaskServer/util/Tools"
	"BTaskServer/util/response"
	"BTaskServer/util/validatorTool"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ctx               = context.Background()
	typeList          = []int{11, 12} // 11关注，12点赞
	currentRequests   int32           // 当前处理中的请求数
	taskOrderCounters = sync.Map{}    // 用于跟踪每个订单的任务分配数
)

type ITaskListController interface {
	GetTaskList(c *gin.Context)      // 获取任务列表（管理员）
	GetTaskListByKey(c *gin.Context) // 搜索列表(管理员)
	GetTask(c *gin.Context)          // 领取任务
	SubmitTask(c *gin.Context)       // 提交任务
}

type TaskListController struct {
	DB                    *gorm.DB
	RDS                   *redis.Client
	MaxRedisTime          int   // redis 最大访问间隔（暂时不再使用）
	TaskTimeOut           int   // 任务提交超时时间
	MaxConcurrentRequests int32 // 最大并发请求数
}

// 1. 修改fastCheck函数，移除时间限制
func fastCheck(redis *redis.Client, account string, orderId int) bool {
	// 只检查账号是否已领取过特定任务，不再限制时间
	key := fmt.Sprintf("task:account:%s:order:%d", account, orderId)
	exists, _ := redis.Exists(ctx, key).Result()
	return exists == 0
}

// 2. 修改GetTask方法，移除全局锁，使用数据库行锁
func (t TaskListController) GetTask(c *gin.Context) {
	// 检查当前请求数是否超出最大并发数
	if atomic.LoadInt32(&currentRequests) >= t.MaxConcurrentRequests {
		response.ServerBad(c, nil, "系统繁忙，请稍后重试")
		return
	}
	atomic.AddInt32(&currentRequests, 1)        // 增加计数
	defer atomic.AddInt32(&currentRequests, -1) // 使用 defer 确保请求结束时减少计数

	var getTaskQuery model.GetTaskQuery
	if !validatorTool.ValidatorQuery[*model.GetTaskQuery](c, &getTaskQuery) {
		return
	}

	// 验证加密
	jccStr := fmt.Sprintf("%s%s%d", getTaskQuery.BlAccount, "jcc", getTaskQuery.Time)
	jccMd5 := Tools.GenMd5(jccStr)
	if getTaskQuery.Jcc != jccMd5 {
		response.Fail(c, nil, "jcc验证失败")
		return
	}

	// 验证type是否正确
	isIn := false
	for _, value := range typeList {
		if value == getTaskQuery.Type {
			isIn = true
			break
		}
	}

	if isIn == false {
		response.Fail(c, nil, "type参数错误,没有此类型的任务")
		return
	}

	// 验证ukey是否存在
	user, b := getCacelUser(getTaskQuery.UserKey)
	if b == false {
		response.ServerBad(c, nil, "ukey验证失败")
		return
	}
	// ukey不能是管理员的
	if user.Authority == 1 {
		response.AuthError(c, nil, "不能用管理员的ukey领取任务")
		return
	}

	// 通过获取的订单类型，查询设置的单价
	priceMap := GetPriceByType(user.AddUserId, getTaskQuery.Type)
	if priceMap["ret"].(bool) == false {
		response.ServerBad(c, nil, "领取失败,单价获取失败")
		return
	}

	price := priceMap["price"].(float64)

	// 检查账号是否已领取过该类型任务
	accountTypeKey := fmt.Sprintf("account:%s:type:%d", getTaskQuery.BlAccount, getTaskQuery.Type)

	// 启动数据库事务
	tx := t.DB.Begin()
	if tx.Error != nil {
		response.ServerBad(c, nil, "领取失败，数据库事务启动失败")
		return
	}

	// 获取当前可用的任务
	whereStr := fmt.Sprintf("UStatus = %d and getTaskType = %d", 0, getTaskQuery.Type)

	// 查询该账号已经完成过的任务ID列表
	var completedOrderIds []int
	if err := tx.Model(&model.TaskLog{}).
		Select("DISTINCT orderId").
		Where("Account = ?", getTaskQuery.BlAccount).
		Pluck("orderId", &completedOrderIds).Error; err != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "领取失败，查询历史任务失败")
		return
	}

	// 转换为map方便查询
	completedOrderMap := make(map[int]bool)
	for _, orderId := range completedOrderIds {
		completedOrderMap[orderId] = true
	}

	// 找出所有可用的任务
	var availableTasks []model.TaskList
	if err := tx.Where(whereStr).
		Where("collectNum < buyNumber").
		Order("id DESC").
		Find(&availableTasks).Error; err != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "领取失败，查询任务失败")
		return
	}

	if len(availableTasks) == 0 {
		tx.Rollback()
		response.Fail(c, nil, "暂无可领取的任务")
		return
	}

	// 尝试找到一个可以分配的任务
	var taskList *model.TaskList
	var allocated int
	var allocatedKey string

	// 遍历所有可用任务
	for i := range availableTasks {
		// 检查是否已经做过这个orderID的任务
		if completedOrderMap[availableTasks[i].OrderId] {
			continue // 跳过已经做过的orderID
		}

		// 使用Redis的SETNX尝试获取任务锁
		taskLockKey := fmt.Sprintf("task:lock:%d", availableTasks[i].ID)
		locked, _ := t.RDS.SetNX(ctx, taskLockKey, "1", 10*time.Second).Result()
		if !locked {
			continue // 如果获取锁失败，尝试下一个任务
		}

		// 成功获取到锁，检查Redis中的计数
		allocatedKey = fmt.Sprintf("task:allocated:%d", availableTasks[i].OrderId)
		allocated, _ = t.RDS.Get(ctx, allocatedKey).Int()

		// 如果Redis中没有记录，从数据库获取
		if allocated == 0 {
			// 再次从数据库获取最新任务状态
			var latestTask model.TaskList
			if err := tx.Set("gorm:query_option", "FOR UPDATE").
				Where("id = ?", availableTasks[i].ID).
				First(&latestTask).Error; err != nil {
				// 释放锁
				t.RDS.Del(ctx, taskLockKey)
				continue
			}

			// 使用数据库中的值
			allocated = latestTask.CollectNum
			// 更新Redis
			t.RDS.Set(ctx, allocatedKey, allocated, 2*time.Hour)

			// 检查是否仍有空间
			if allocated >= latestTask.BuyNumber {
				// 释放锁
				t.RDS.Del(ctx, taskLockKey)
				continue
			}

			taskList = &latestTask
		} else {
			// 检查Redis中的计数是否已满
			if allocated >= availableTasks[i].BuyNumber {
				// 释放锁
				t.RDS.Del(ctx, taskLockKey)
				continue
			}

			// 再次从数据库确认最新状态
			var latestTask model.TaskList
			if err := tx.Set("gorm:query_option", "FOR UPDATE").
				Where("id = ?", availableTasks[i].ID).
				First(&latestTask).Error; err != nil {
				// 释放锁
				t.RDS.Del(ctx, taskLockKey)
				continue
			}

			// 最终确认是否有空间
			if latestTask.CollectNum >= latestTask.BuyNumber {
				// 释放锁
				t.RDS.Del(ctx, taskLockKey)
				continue
			}

			taskList = &latestTask
		}

		// 找到了可用任务，跳出循环
		break
	}

	// 如果没有找到可用任务
	if taskList == nil {
		tx.Rollback()
		response.Fail(c, nil, "暂无可领取的任务")
		return
	}

	// 获取任务锁名称并确保在函数结束时释放
	taskLockKey := fmt.Sprintf("task:lock:%d", taskList.ID)
	defer t.RDS.Del(ctx, taskLockKey)

	// 增加Redis中的分配计数
	newAllocated, err := t.RDS.Incr(ctx, allocatedKey).Result()
	if err != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "领取失败，Redis计数失败")
		return
	}
	t.RDS.Expire(ctx, allocatedKey, 2*time.Hour)

	// 最终确认分配数量不超过限制
	if newAllocated > int64(taskList.BuyNumber) {
		// 如果超过了限制，回滚Redis计数
		t.RDS.Decr(ctx, allocatedKey)
		tx.Rollback()
		response.Fail(c, nil, "任务已被全部领取")
		return
	}

	// 设置账号任务关联
	accountTaskKey := fmt.Sprintf("account:%s:order:%d", getTaskQuery.BlAccount, taskList.OrderId)
	err = t.RDS.SetNX(ctx, accountTaskKey, 1, 6*time.Hour).Err()
	if err != nil {
		// 回滚Redis计数
		t.RDS.Decr(ctx, allocatedKey)
		tx.Rollback()
		response.ServerBad(c, nil, "领取失败，Redis操作失败")
		return
	}

	// 创建任务日志
	newTaskLog := model.TaskLog{
		CollectTime: Tools.GetDateNowFormat(true),
		OrderNum:    Tools.GetUuid(),
		Account:     getTaskQuery.BlAccount,
		VideoLink:   taskList.VideoLink,
		OrderId:     taskList.OrderId,
		GoodsId:     taskList.GoodsId,
		GoodsName:   taskList.GoodsName,
		Price:       price,
		UserKey:     user.UserKey,
		ManagerId:   user.AddUserId,
		Status:      0,
		GetTaskType: taskList.GetTaskType,
		PingtaiName: taskList.PingtaiName,
	}

	// 执行最终安全检查，确保数据库中的领取数量不会超过限制
	// 这是一个关键的双重保险机制，确保在高并发情况下仍能防止任务超领
	var totalTaskCount int64
	if err := tx.Model(&model.TaskList{}).
		Where("orderId = ?", taskList.OrderId).
		Select("collectNum").
		Pluck("collectNum", &totalTaskCount).Error; err != nil {
		// 查询失败，回滚Redis计数
		t.RDS.Decr(ctx, allocatedKey)
		t.RDS.Del(ctx, accountTaskKey)
		t.RDS.Del(ctx, accountTypeKey)
		tx.Rollback()
		response.ServerBad(c, nil, "领取失败，无法获取任务总计数")
		return
	}

	// 如果当前数据库中的计数+1已经超过限制，则不允许领取
	// 这是防止任务超领的最后一道防线，直接基于数据库当前状态判断，避免Redis计数不准确的问题
	if totalTaskCount+1 > int64(taskList.BuyNumber) {
		// 回滚Redis计数
		t.RDS.Decr(ctx, allocatedKey)
		t.RDS.Del(ctx, accountTaskKey)
		t.RDS.Del(ctx, accountTypeKey)
		tx.Rollback()
		response.Fail(c, nil, "任务领取数量已达上限")
		return
	}

	// 更新数据库任务领取数量
	if err := tx.Model(&model.TaskList{}).Where("id = ? AND collectNum < buyNumber", taskList.ID).
		Updates(map[string]interface{}{
			"collectNum":  gorm.Expr("collectNum + 1"),
			"lastGetTime": Tools.GetDateNowFormat(true),
		}).Error; err != nil {
		// 回滚所有Redis设置
		t.RDS.Del(ctx, accountTaskKey)
		t.RDS.Del(ctx, accountTypeKey)
		t.RDS.Decr(ctx, allocatedKey)
		tx.Rollback()
		response.ServerBad(c, nil, "领取失败，更新任务数据失败")
		return
	}

	// 验证更新结果
	var updatedTask model.TaskList
	if err := tx.Where("id = ?", taskList.ID).First(&updatedTask).Error; err != nil {
		// 回滚所有Redis设置
		t.RDS.Del(ctx, accountTaskKey)
		t.RDS.Del(ctx, accountTypeKey)
		t.RDS.Decr(ctx, allocatedKey)
		tx.Rollback()
		response.ServerBad(c, nil, "领取失败，验证任务数据失败")
		return
	}

	// 再次严格检查数据库中的数量是否超限
	if updatedTask.CollectNum > updatedTask.BuyNumber {
		// 回滚所有Redis设置
		t.RDS.Del(ctx, accountTaskKey)
		t.RDS.Del(ctx, accountTypeKey)
		t.RDS.Decr(ctx, allocatedKey)
		tx.Rollback()
		response.Fail(c, nil, "任务领取数量已达上限")
		return
	}

	// 创建任务日志
	if err := tx.Create(&newTaskLog).Error; err != nil {
		// 回滚所有Redis设置
		t.RDS.Del(ctx, accountTaskKey)
		t.RDS.Del(ctx, accountTypeKey)
		t.RDS.Decr(ctx, allocatedKey)
		tx.Rollback()
		response.ServerBad(c, nil, "领取失败，创建任务日志失败")
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		// 回滚所有Redis设置
		t.RDS.Del(ctx, accountTaskKey)
		t.RDS.Del(ctx, accountTypeKey)
		t.RDS.Decr(ctx, allocatedKey)
		tx.Rollback()
		response.ServerBad(c, nil, "领取失败，提交事务失败")
		return
	}

	// 同步总任务数到Redis（用于提交阶段验证）
	taskTotalKey := fmt.Sprintf("task:total:%d", taskList.OrderId)
	t.RDS.Set(ctx, taskTotalKey, taskList.BuyNumber, 6*time.Hour)

	resultTask := model.GetTaskModel{
		ID:        int(newTaskLog.ID),
		VideoLink: taskList.VideoLink,
		ShortLink: taskList.ShortLink,
		VideoId:   taskList.VideoId,
		Uid:       taskList.Uid,
		GoodsName: taskList.GoodsName,
	}

	fmt.Println(fmt.Sprintf("%s,orderid:%d,ks号:%s,%s", taskList.PingtaiName, taskList.OrderId, getTaskQuery.BlAccount, Tools.GetDateNowFormat(true)))

	response.Success(c, gin.H{"task": resultTask}, fmt.Sprintf("任务领取成功,请在%d分钟内提交", t.TaskTimeOut))
}

// 回滚Redis分配计数的辅助函数
func rollbackRedisAllocation(rds *redis.Client, allocatedKey, accountKey string) {
	pipe := rds.Pipeline()
	pipe.Decr(ctx, allocatedKey)
	pipe.Del(ctx, accountKey)
	pipe.Exec(ctx)
}

// 3. 修改SubmitTask方法，使用新的键名方案
func (t TaskListController) SubmitTask(c *gin.Context) {
	var submitTaskQuery model.SubmitTaskQuery
	if !validatorTool.ValidatorQuery[*model.SubmitTaskQuery](c, &submitTaskQuery) {
		return
	}

	// 启动数据库事务
	tx := t.DB.Begin()
	if tx.Error != nil {
		response.ServerBad(c, nil, "提交失败，事务启动失败")
		return
	}

	// 获取任务日志并锁定
	var taskLog model.TaskLog
	res := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", submitTaskQuery.ID).
		Limit(1).Find(&taskLog)

	if res.Error != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "提交失败，查询任务日志失败")
		return
	}

	if res.RowsAffected == 0 {
		tx.Rollback()
		response.Fail(c, nil, "提交失败，未找到领取记录")
		return
	}

	if taskLog.Status == 2 {
		tx.Rollback()
		response.Fail(c, nil, "提交失败，任务已超时")
		return
	}

	if taskLog.Status != 0 {
		tx.Rollback()
		response.Fail(c, nil, "提交失败，不能重复提交")
		return
	}

	// 获取任务信息
	var taskInfo model.TaskList
	if err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("orderId = ?", taskLog.OrderId).First(&taskInfo).Error; err != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "提交失败，任务信息不存在")
		return
	}

	// 更新任务日志状态为已提交
	if err := updateTaskLogSuccess(tx, taskLog.ID); err != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "提交失败，更新任务日志失败")
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.ServerBad(c, nil, "提交失败，事务提交失败")
		return
	}

	response.Success(c, nil, "提交成功")
}

// 辅助函数：更新任务日志状态为成功
func updateTaskLogSuccess(db *gorm.DB, taskLogId uint) error {
	updates := map[string]interface{}{
		"status":      1, // 已提交等待审核
		"sumbmitTime": Tools.GetDateNowFormat(true),
	}

	return db.Model(&model.TaskLog{}).Where("id = ? AND status = 0", taskLogId).
		Updates(updates).Error
}

// 4. 修改NewTaskListController构造函数，增加并发支持
func NewTaskListController() ITaskListController {
	db := common.GetDB()
	rds := common.GetRedis()

	if err := db.AutoMigrate(&model.TaskList{}); err != nil {
		panic("taskList表迁移失败")
	}
	fmt.Println("taskList表迁移成功")

	timeout := viper.GetInt("TaskConfig.TaskTimeOut")
	maxConcurrentRequests := viper.GetInt32("system.maxConcurrentRequests")

	// 初始化控制器，设置高并发请求数
	return TaskListController{
		DB:                    db,
		RDS:                   rds,
		TaskTimeOut:           timeout,               // 任务超时时间（分钟）
		MaxConcurrentRequests: maxConcurrentRequests, // 设置高并发请求上限
	}
}

func (t TaskListController) GetTaskListByKey(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	var findTaskQuery model.FindTaskQuery
	if !validatorTool.ValidatorQuery[*model.FindTaskQuery](c, &findTaskQuery) {
		return
	}

	orderStr := "%" + findTaskQuery.OrderId + "%"

	var taskList []model.TaskList
	if res := t.DB.Raw("select a.* from taskList as a join (select id from taskList where orderId like ? order by orderId limit ?,?) b on a.id = b.id", orderStr, ((*findTaskQuery.PageNum - 1) * *findTaskQuery.PageSize), *findTaskQuery.PageSize).Scan(&taskList); res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	t.DB.Model(model.TaskList{}).Where("orderId like ?", orderStr).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *findTaskQuery.PageNum,
		"result":  taskList,
	}, "获取成功")
}

func (t TaskListController) GetTaskList(c *gin.Context) {
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
	var taskList []model.TaskList
	if res := t.DB.Raw("select a.* from taskList as a join (select id from taskList order by downTime desc limit ?,?) b on a.id = b.id", ((*pageQuery.PageNum - 1) * *pageQuery.PageSize), *pageQuery.PageSize).Scan(&taskList); res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	t.DB.Model(model.TaskList{}).Where("1=1").Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *pageQuery.PageNum,
		"result":  taskList,
	}, "获取成功")
}
