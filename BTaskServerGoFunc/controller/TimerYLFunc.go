package controller

import (
	"BTaskServer/model"
	"BTaskServer/util/BLTaskFunc"
	"BTaskServer/util/Tools"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"strings"
	"time"
)

// 亿乐
//1	已付款
//2	待处理
//3	处理中
//4	补单中
//5	退单中
//6	已完成
//7	已退单
//8	已退款
//9	有异常

var page int
var SchedulePage int
var typeList []int = []int{3511, 3512} // 亿乐type

// 亿乐拉单
func GetHigherOrdersYL(db *gorm.DB) {
	timer := viper.GetInt("Timer.DownOrderTimer")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("亿乐拉取订单recover%s", Tools.GetDateNowFormat(true)))
			go GetHigherOrdersYL(db)
		}
	}()

	LDpage := 1

	for {
		fmt.Println(fmt.Sprintf("亿乐拉取订单%s,page%d", Tools.GetDateNowFormat(true), LDpage))
		resMap := BLTaskFunc.GetOrderListYL(LDpage, 50, 3)
		if resMap["ret"].(bool) == false {
			fmt.Println("亿乐拉取订单失败，" + resMap["msg"].(string))
			time.Sleep(time.Second)
			continue
		}

		data := resMap["data"].(map[string]interface{})
		infos := data["infos"].([]interface{})

		if len(infos) == 0 {
			fmt.Println(fmt.Sprintf("亿乐拉取订单完成len=0 %s", Tools.GetDateNowFormat(true)))
			time.Sleep(time.Second * time.Duration(timer))
			LDpage = 1
			continue
		}

		LDpage++

		var orderIds []int
		var taskListS []model.TaskList

		for _, value := range infos {
			if value == nil {
				continue
			}

			valueObj := value.(map[string]interface{})
			goods_id := int(valueObj["goods_id"].(float64))

			if isCodeIn(typeList, goods_id) == false {
				continue
			}
			orderId := int(valueObj["id"].(float64))
			buyNum := int(valueObj["buy_number"].(float64))
			shortId := valueObj["parameter"].(string)     // 短链接
			videoLongLink, _ := Tools.GetLongUrl(shortId) // 长链接

			//如果长链接中包含 “/user” 说明？之前的是uid
			uid := ""
			videoId := ""
			if strings.Index(videoLongLink, "/user") != -1 {
				uid = analysisKSVid(videoLongLink) // uid
			} else {
				videoId = analysisKSVid(videoLongLink)          // 视频id
				uid = analysisKSAnyId(videoLongLink, "userId=") //uid
			}

			var getTaskType int
			if goods_id == 3511 {
				getTaskType = 11
			} else if goods_id == 3512 {
				getTaskType = 12
			}

			insertMap := model.TaskList{
				DownTime:     Tools.GetDateNowFormat(true),
				VideoLink:    videoLongLink,
				ShortLink:    shortId,
				VideoId:      videoId,
				Uid:          uid,
				OrderId:      orderId,
				GoodsId:      goods_id,
				SellingPrice: valueObj["selling_price"].(float64),
				BuyNumber:    buyNum,
				Amount:       valueObj["amount"].(float64),
				StartNum:     int(valueObj["start_num"].(float64)),
				CurrentNum:   0,
				CollectNum:   0,
				GoodsStatus:  3,
				Remark:       valueObj["remark"].(string),
				GoodsName:    valueObj["goods_name"].(string),
				UStatus:      0,
				GetTaskType:  getTaskType,
				PingtaiName:  "亿乐",
			}

			orderIds = append(orderIds, orderId)
			taskListS = append(taskListS, insertMap)
		}

		var existingRecords []model.TaskList
		db.Where("orderId in ?", orderIds).Find(&existingRecords)

		// 构建已存在的 集合
		existMap := make(map[int]bool)
		for _, r := range existingRecords {
			existMap[r.OrderId] = true
		}

		// 批量插入不存在的数据
		var newRecords []model.TaskList
		for _, orderid := range orderIds {
			if !existMap[orderid] {
				for _, tasklist := range taskListS {
					if orderid == tasklist.OrderId {
						newRecords = append(newRecords, tasklist)
					}
				}

			}
		}

		fmt.Println(fmt.Sprintf("亿乐插入新订单 %d 条 %s", len(newRecords), Tools.GetDateNowFormat(true)))

		if len(newRecords) > 0 {
			db.Create(&newRecords)
		}

		time.Sleep(time.Second * time.Duration(timer))
	}

}

// 亿乐退单
func GetTDOrdersYL(db *gorm.DB, status int) {
	timer := viper.GetInt("Timer.UpdateTDOrderTimer")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("亿乐退单recover%s", Tools.GetDateNowFormat(true)))
			go GetTDOrdersYL(db, status)
		}
	}()

	TDpage := 1

	for {
		tx := db.Begin()
		resMap := BLTaskFunc.GetOrderListYL(TDpage, 50, status)
		fmt.Println(fmt.Sprintf("亿乐退单%d订单%s,page%d", status, Tools.GetDateNowFormat(true), TDpage))
		if resMap["ret"].(bool) == false {
			fmt.Println(fmt.Sprintf("亿乐退单%d订单失败,%s", status, resMap["msg"].(string)))
			time.Sleep(time.Second)
			continue
		}
		data := resMap["data"].(map[string]interface{})
		infos := data["infos"].([]interface{})

		if len(infos) == 0 {
			time.Sleep(time.Second * time.Duration(timer))
			TDpage = 1
			continue
		}

		TDpage++

		var orderIdList []int
		for _, value := range infos {
			if value == nil {
				continue
			}

			valueObj := value.(map[string]interface{})

			orderid := int(valueObj["id"].(float64))
			orderIdList = append(orderIdList, orderid)
		}

		//db.Model(&model.TaskList{}).Where("orderId in ?", orderIdList).Update("UStatus", 2)
		if len(orderIdList) > 0 {
			res := tx.Model(&model.TaskList{}).Where("orderId in ? and UStatus <> 2", orderIdList).Updates(map[string]interface{}{
				"UStatus":     2,
				"goodsStatus": status,
				"remark":      fmt.Sprintf("该任务已退单,%s", Tools.GetDateNowFormat(true)),
			})
			if res.Error != nil {
				tx.Rollback()
				time.Sleep(time.Second * time.Duration(timer))
				continue
			}

			if res.RowsAffected == 0 {
				tx.Rollback()
				time.Sleep(time.Second * time.Duration(timer))
				continue
			}

			tx.Commit()
		}

		fmt.Println(fmt.Sprintf("亿乐退单%d订单%s完成", status, Tools.GetDateNowFormat(true)))
		time.Sleep(time.Second * time.Duration(timer))
	}
}

// 拉取已完成的订单，更新本地的状态
// 亿乐完成
func GetWCOrdersYL(db *gorm.DB, status int) {
	timer := viper.GetInt("Timer.UpdateTDOrderTimer")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("亿乐退单recover%s", Tools.GetDateNowFormat(true)))
			go GetTDOrdersYL(db, status)
		}
	}()

	TDpage := 1

	for {
		tx := db.Begin()
		resMap := BLTaskFunc.GetOrderListYL(TDpage, 50, status)
		fmt.Println(fmt.Sprintf("亿乐退单%d订单%s,page%d", status, Tools.GetDateNowFormat(true), TDpage))
		if resMap["ret"].(bool) == false {
			fmt.Println(fmt.Sprintf("亿乐退单%d订单失败,%s", status, resMap["msg"].(string)))
			time.Sleep(time.Second)
			continue
		}
		data := resMap["data"].(map[string]interface{})
		infos := data["infos"].([]interface{})

		if len(infos) == 0 {
			time.Sleep(time.Second * time.Duration(timer))
			TDpage = 1
			continue
		}

		TDpage++

		var orderIdList []int
		for _, value := range infos {
			if value == nil {
				continue
			}

			valueObj := value.(map[string]interface{})

			orderid := int(valueObj["id"].(float64))
			orderIdList = append(orderIdList, orderid)
		}

		//db.Model(&model.TaskList{}).Where("orderId in ?", orderIdList).Update("UStatus", 2)
		if len(orderIdList) > 0 {
			res := tx.Model(&model.TaskList{}).Where("orderId in ? and UStatus <> 2", orderIdList).Updates(map[string]interface{}{
				"UStatus":     2,
				"goodsStatus": status,
				"remark":      fmt.Sprintf("该任务已退单,%s", Tools.GetDateNowFormat(true)),
			})
			if res.Error != nil {
				tx.Rollback()
				time.Sleep(time.Second * time.Duration(timer))
				continue
			}

			if res.RowsAffected == 0 {
				tx.Rollback()
				time.Sleep(time.Second * time.Duration(timer))
				continue
			}

			tx.Commit()
		}

		fmt.Println(fmt.Sprintf("亿乐退单%d订单%s完成", status, Tools.GetDateNowFormat(true)))
		time.Sleep(time.Second * time.Duration(timer))
	}
}

// 定时更新供货端进度
func OrderScheduleYL(db *gorm.DB) {
	go func() {

		SchedulePage = 1
		for {
			fmt.Println(fmt.Sprintf("更新进度%d", SchedulePage))

			var taskList []model.TaskList

			res := db.Raw("select a.* from taskList as a join (select id from taskList WHERE UStatus = 0 ORDER BY downTime desc limit ?,?) b on a.id = b.id", (SchedulePage-1)*50, 50).Scan(&taskList)
			if res.Error != nil {
				time.Sleep(time.Second * 1)
				continue
			}

			if res.RowsAffected == 0 {
				SchedulePage = 1
				time.Sleep(time.Second * 60)
				continue
			}

			for _, task := range taskList {
				BLTaskFunc.UpdateOrderScheduleYL(task.OrderId, task.StartNum, task.CurrentNum)
				time.Sleep(time.Millisecond * 100)
			}

			SchedulePage++
			time.Sleep(time.Second * 1)
			fmt.Println(fmt.Sprintf("更新进度完成%d", SchedulePage))
		}

	}()

}
