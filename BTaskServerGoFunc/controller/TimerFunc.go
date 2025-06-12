package controller

import (
	"BTaskServer/model"
	"BTaskServer/util/BLTaskFunc"
	"BTaskServer/util/Tools"
	"BTaskServer/util/WorkPool"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

var ctx = context.Background()

func Shenghe(db *gorm.DB) {

	timer := viper.GetInt("Timer.ShengheTimer")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("执行审核recover%s", Tools.GetDateNowFormat(true)))
			go Shenghe(db)
		}
	}()

	for {
		fmt.Println(fmt.Sprintf("执行审核%s", Tools.GetDateNowFormat(true)))
		manager := getCacelGaiTc()
		ShTcGl := manager.ShTcGl

		var userList []model.User
		res1 := db.Where("authority = ?", 0).Find(&userList)
		if res1.Error != nil {
			time.Sleep(time.Second * time.Duration(timer))
			continue
		}

		if res1.RowsAffected == 0 {
			time.Sleep(time.Second * time.Duration(timer))
			continue
		}

		pool := WorkPool.NewWorkerPool(10)
		for _, user := range userList {
			userModel := user
			pool.Submit(func() {
				ShengheThread(db, ShTcGl, userModel)
			})
		}
		pool.Wait()
		fmt.Println(fmt.Sprintf("本轮执行审核完成%s", Tools.GetDateNowFormat(true)))
		time.Sleep(time.Second * time.Duration(timer))
	}
}

func ShengheThread(db *gorm.DB, ShTcGl int, user model.User) {
	var addMoney float64
	var sucessIds []uint
	var failIds []uint

	var taskLogList []model.TaskLog
	res2 := db.Where("userKey = ? and status = 1", user.UserKey).Limit(500).Find(&taskLogList)
	if res2.Error != nil {
		return
	}

	// 待审大于300条才执行
	//if res2.RowsAffected < 300 {
	//	return
	//}

	if res2.RowsAffected == 0 {
		return
	}

	//var diggCount int
	//var fenCount int
	for _, taskLog := range taskLogList {
		r1 := rand.Intn(100)
		if r1 >= ShTcGl {
			// 不偷吃
			addMoney = addMoney + taskLog.Price
			sucessIds = append(sucessIds, taskLog.ID)
		} else {
			// 偷吃
			failIds = append(failIds, taskLog.ID)
		}
	}

	//newShengheLog := model.ShengheLog{
	//	UserId:         user.ID,
	//	UserKey:        user.UserKey,
	//	AllTotal:       len(taskLogList),
	//	EffectiveTotal: len(sucessIds),
	//	DiggTotal:      diggCount,
	//	FenTotal:       fenCount,
	//	AddPrice:       addMoney,
	//	ExamineTime:    Tools.GetDateNowFormat(true),
	//}
	//
	//if errSH := db.Create(&newShengheLog); errSH != nil {
	//	return
	//}

	tx := db.Begin()

	// 不偷吃
	res10 := tx.Model(&model.TaskLog{}).Where("id in ?", sucessIds).Updates(map[string]interface{}{
		"status":      3,
		"examineTime": Tools.GetDateNowFormat(true),
	})
	if res10.Error != nil {
		tx.Rollback()
		return
	}

	if res10.RowsAffected == 0 {
		tx.Rollback()
		return
	}

	//更新user的money
	res3 := tx.Model(&model.User{}).Where("userKey = ?", user.UserKey).Update("money", gorm.Expr(fmt.Sprintf("money + %f", addMoney)))
	if res3.Error != nil {
		tx.Rollback()
		return
	}

	if res3.RowsAffected == 0 {
		tx.Rollback()
		return
	}

	// 偷吃
	res11 := tx.Model(&model.TaskLog{}).Where("id in ?", failIds).Updates(map[string]interface{}{
		"status":      4,
		"examineTime": Tools.GetDateNowFormat(true),
	})
	if res11.Error != nil {
		tx.Rollback()
		return
	}

	if res11.RowsAffected == 0 {
		tx.Rollback()
		return
	}

	tx.Commit()

	fmt.Println(fmt.Sprintf("审核 用户id:%d %d条 增加金额%f", user.ID, len(taskLogList), addMoney))
}

// 更新完成的状态，不修改上级，通用
func SubHigherOrders(db *gorm.DB) {
	timer := viper.GetInt("Timer.UpOrderTimer")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("更改完成订单recover%s", Tools.GetDateNowFormat(true)))
			go SubHigherOrders(db)
		}
	}()

	for {
		fmt.Println("更改完成订单" + Tools.GetDateNowFormat(true))
		//tx := db.Begin()
		//var TaskLists []model.TaskList
		//// 播放任务，状态改成3
		//res := db.Where("currentNum >= buyNumber and UStatus = 0").Limit(30).Find(&TaskLists).Updates(map[string]interface{}{
		//	"currentNum":  gorm.Expr("buyNumber"),
		//	"collectNum":  gorm.Expr("buyNumber"),
		//	"UStatus":     1,
		//	"doneTime":    Tools.GetDateNowFormat(true),
		//	"goodsStatus": 6,
		//})

		res := db.Model(&model.TaskList{}).Where("currentNum >= buyNumber and UStatus = 0").Limit(30).Updates(map[string]interface{}{
			"currentNum":  gorm.Expr("buyNumber"),
			"collectNum":  gorm.Expr("buyNumber"),
			"UStatus":     1,
			"doneTime":    Tools.GetDateNowFormat(true),
			"goodsStatus": 6,
		})

		if res.Error != nil {
			time.Sleep(time.Second * time.Duration(timer))
			continue
		}
		if res.RowsAffected == 0 {
			time.Sleep(time.Second * time.Duration(timer))
			continue
		}

		////ret := true
		//for _, taskList := range TaskLists {
		//
		//	// 更新这个订单的进度
		//	retMap1 := BLTaskFunc.UpdateOrderScheduleYL(taskList.OrderId, taskList.StartNum, taskList.BuyNumber)
		//	// 更新这个订单的状态为完成
		//	retMap2 := BLTaskFunc.UpdateOrderStatusYL(taskList.OrderId, 3, 6)
		//
		//	if retMap1["ret"].(bool) == false || retMap2["ret"].(bool) == false {
		//		fmt.Println(retMap1["msg"].(string))
		//		fmt.Println(retMap2["msg"].(string))
		//		//ret = false
		//		//break
		//	}
		//
		//	time.Sleep(time.Second * time.Duration(timer))
		//}
		fmt.Println("更改完成订单完成" + Tools.GetDateNowFormat(true))
	}
}

func DelData(db *gorm.DB) {
	timeout := viper.GetString("TaskConfig.DelTimeOut")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("定时清除数据recover%s", Tools.GetDateNowFormat(true)))
			go DelData(db)
		}
	}()

	for {
		fmt.Println(fmt.Sprintf("定时清除数据%s", Tools.GetDateNowFormat(true)))

		//db.Where("UStatus <> 0 and lastGetTime <> '' and lastGetTime < DATE_SUB(NOW(),INTERVAL ? day)", timeout).Delete(&model.TaskList{})
		db.Where("UStatus <> 2 and downTime < DATE_SUB(NOW(),INTERVAL ? day)", timeout).Delete(&model.TaskList{})
		time.Sleep(time.Second * 2)
		db.Where("UStatus = ?", 2).Delete(&model.TaskList{})

		time.Sleep(time.Second * 2)
		db.Where("orderId not in (select orderId from taskList)").Delete(&model.TaskLog{})
		fmt.Println(fmt.Sprintf("定时清除数据完成%s", Tools.GetDateNowFormat(true)))
		time.Sleep(time.Second * 90)
	}

}

// 开启一个goroutine每小时更新一次，领取记录时间超过x分钟的，修改状态，修改对应的数量
func UpdateTaskNoGet(db *gorm.DB, rds *redis.Client) {

	timeout := viper.GetString("TaskConfig.TaskTimeOut")
	//clearTime := viper.GetString("TaskConfig.CLTimerH")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("清除未提交recover%s", Tools.GetDateNowFormat(true)))
			go UpdateTaskNoGet(db, rds)
		}
	}()

	for {
		fmt.Println(fmt.Sprintf("清除未提交%s", Tools.GetDateNowFormat(true)))
		res := db.Model(&model.TaskLog{}).Where("status = 0 and collectTime <= DATE_SUB(NOW(),INTERVAL ? MINUTE)", timeout).Updates(map[string]interface{}{
			"status": 2,
			"remark": fmt.Sprintf("超过%s分钟未提交", timeout),
		})

		if res.RowsAffected == 0 {
			time.Sleep(time.Second * 2)
			continue
		}

		// 每隔六小时清除一次未提交
		//db.Where("status = 2 and collectTime <= DATE_SUB(NOW(),INTERVAL ? HOUR)", clearTime).Delete(&model.TaskLog{})

		var taskList []model.TaskList
		res1 := db.Where("UStatus = 0").Find(&taskList)
		if res1.Error != nil {
			time.Sleep(time.Second * 2)
			continue
		}

		pool := WorkPool.NewWorkerPool(5)
		for _, taskListObj := range taskList {
			taskListModel := taskListObj
			pool.Submit(func() {
				UpdateTaskNoGetThread(db, rds, taskListModel)
			})
		}
		pool.Wait()
		fmt.Println(fmt.Sprintf("清除未提交完成%s", Tools.GetDateNowFormat(true)))
		time.Sleep(time.Second * 70) // 每70秒清理一次
	}

}

func UpdateTaskNoGetThread(db *gorm.DB, rds *redis.Client, taskListObj model.TaskList) {
	var count int64 // 领取数量
	res2 := db.Model(&model.TaskLog{}).Where("orderId = ? and status in (0,3,4)", taskListObj.OrderId).Count(&count)
	if res2.Error != nil {
		return
	}

	var count2 int64 // 提交数量
	res3 := db.Model(&model.TaskLog{}).Where("orderId = ? and status in (3,4)", taskListObj.OrderId).Count(&count2)
	if res3.Error != nil {
		return
	}

	taskListObj.CurrentNum = int(count2)
	taskListObj.CollectNum = int(count)

	// 修改redis缓存的数量
	fmt.Println(fmt.Sprintf("修改order:%d已领取为:%d", taskListObj.OrderId, count))
	allocatedKey := fmt.Sprintf("task:allocated:%d", taskListObj.OrderId)
	rds.Set(ctx, allocatedKey, count, 2*time.Hour)

	db.Save(&taskListObj)
}

func UpdateDoneList(db *gorm.DB) {
	timer := viper.GetInt("Timer.UpdateDoneTimer")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("更新完成状态recover%s", Tools.GetDateNowFormat(true)))
			go UpdateDoneList(db)
		}
	}()

	for {
		var taskList []model.TaskList
		res1 := db.Where("UStatus in (0,1) and currentNum >= buyNumber and videoId = ''").Find(&taskList)
		if res1.Error != nil {
			time.Sleep(time.Second * 2)
			continue
		}

		var status int
		for _, taskListObj := range taskList {
			var count int64
			res2 := db.Model(&model.TaskLog{}).Where("orderId = ? and status in (3,4)", taskListObj.OrderId).Count(&count)
			if res2.Error != nil {
				continue
			}

			if int(count) >= taskListObj.BuyNumber {
				status = 1
			} else {
				status = 0
			}

			//db.Model(&model.TaskList{}).Where("id = ?", taskListObj.ID).Updates(map[string]interface{}{
			//	"currentNum": count,
			//	"UStatus":    status,
			//})

			taskListObj.CurrentNum = int(count)
			taskListObj.UStatus = status
			db.Save(&taskListObj)
		}

		fmt.Println(fmt.Sprintf("更新完成状态完成%s", Tools.GetDateNowFormat(true)))
		time.Sleep(time.Second * time.Duration(timer)) // 每21秒清理一次
	}
}

func UpdateDoneList01(db *gorm.DB) {
	timer := viper.GetInt("Timer.UpdateDoneTimer")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("更新完成状态recover%s", Tools.GetDateNowFormat(true)))
			go UpdateDoneList01(db)
		}
	}()

	for {
		var taskList []model.TaskList
		res1 := db.Where("UStatus in (0,1) and currentNum >= buyNumber and errCount < 2").Find(&taskList)
		if res1.Error != nil {
			time.Sleep(time.Second * 2)
			continue
		}

		pool := WorkPool.NewWorkerPool(5)
		for _, taskListObj := range taskList {
			taskListModel := taskListObj
			pool.Submit(func() {
				UpdateDoneList01Thread(db, taskListModel)
			})
		}
		pool.Wait()

		fmt.Println(fmt.Sprintf("更新完成状态完成%s", Tools.GetDateNowFormat(true)))
		time.Sleep(time.Second * time.Duration(timer)) // 每21秒清理一次
	}
}

func UpdateDoneList01Thread(db *gorm.DB, taskListModel model.TaskList) {

	resMap := BLTaskFunc.GetKSCount(taskListModel.ShortLink, taskListModel.GetTaskType)
	if resMap["ret"].(bool) == false {
		fmt.Println(fmt.Sprintf("id %d,错误 %s", taskListModel.ID, taskListModel.ShortLink))
		taskListModel.ErrCount = taskListModel.ErrCount + 1
		db.Save(&taskListModel)
		return
	}

	//var nowCount int // 现在的实际数量
	//if taskListModel.GetTaskType == 11 {
	//	fans := resMap["count"].(int)
	//	nowCount = fans - taskListModel.StartNum
	//} else if taskListModel.GetTaskType == 12 {
	//	like := resMap["count"].(int)
	//	nowCount = like - taskListModel.StartNum
	//}
	count := resMap["count"].(int)
	nowCount := count - taskListModel.StartNum

	if nowCount < 0 {
		return
	}

	if nowCount >= taskListModel.BuyNumber {
		//db.Model(&taskListModel).Update("UStatus", 4)
		fmt.Println(fmt.Sprintf("id %d,实际数量:%d,需要数量%d,完成", taskListModel.ID, nowCount, taskListModel.BuyNumber))
		taskListModel.UStatus = 4
		taskListModel.CurrentNum = nowCount
		taskListModel.CollectNum = nowCount
		db.Save(&taskListModel)
		return
	}

	gaiCount := nowCount

	fmt.Println(fmt.Sprintf("id %d,实际数量:%d,需要数量%d,改%d", taskListModel.ID, nowCount, taskListModel.BuyNumber, gaiCount))

	taskListModel.UStatus = 0
	taskListModel.CurrentNum = gaiCount
	taskListModel.CollectNum = gaiCount
	db.Save(&taskListModel)
}

func isCodeIn(codeList []int, code int) bool {
	ret := false
	for _, codeTemp := range codeList {
		if codeTemp == code {
			ret = true
			break
		}
	}
	return ret
}

func analysisKSVid(videoLongLink string) string {
	videoKey1 := "?" // 视频id关键字1
	index1 := strings.Index(videoLongLink, videoKey1)
	if index1 == -1 {
		return ""
	}

	linkList1 := strings.Split(videoLongLink, videoKey1)
	link1 := linkList1[0]

	videoKey2 := "/"

	linkList2 := strings.Split(link1, videoKey2)
	link2 := linkList2[len(linkList2)-1]

	return link2
}

func analysisKSAnyId(videoLongLink string, keys string) string {
	index1 := strings.Index(videoLongLink, keys)
	if index1 == -1 {
		return ""
	}

	List1 := strings.Split(videoLongLink, "&")
	shortLink := ""
	for _, v := range List1 {
		index11 := strings.Index(v, keys)
		if index11 != -1 {
			shortLink = strings.ReplaceAll(v, keys, "")
			break
		}
	}

	return shortLink
}

func analysisBL(videoLongLink string) string {

	// 包含b23.tv的，需要短链接转长链接
	indextv := strings.Index(videoLongLink, "b23.tv")
	if indextv != -1 {
		videoLongLink, _ = Tools.GetLongUrl(videoLongLink)
	}

	// 处理带video的
	indexStr := "video/"
	indexvideo := strings.Index(videoLongLink, indexStr)

	if indexvideo == -1 {
		indexStr = "bilibili.com/"
		indexvideo = strings.Index(videoLongLink, indexStr)
	}

	indexfoot := strings.Index(videoLongLink, "?")
	if indexfoot == -1 {
		indexfoot = len(videoLongLink)
	}

	shortLink := ""
	shortLink = videoLongLink[indexvideo+len(indexStr) : indexfoot]
	shortLink = strings.ReplaceAll(shortLink, "/", "")

	return shortLink
}

func addOrder(db *gorm.DB, insertMap model.TaskList) {
	var taskList model.TaskList
	res := db.Where("orderId = ?", insertMap.OrderId).Limit(1).Find(&taskList)
	if res.Error != nil {
		return
	}

	if res.RowsAffected != 0 {
		return
	}

	tx := db.Begin()

	if res1 := tx.Create(&insertMap); res1.Error != nil {
		tx.Rollback()
		return
	}

	////更新这个订单的状态为待处理
	//resMap := BLTaskFunc.UpdateOrderStatusYL(insertMap.OrderId, 1, 2)
	//if resMap["ret"].(bool) == false {
	//	tx.Rollback()
	//	return
	//}

	tx.Commit()
}
