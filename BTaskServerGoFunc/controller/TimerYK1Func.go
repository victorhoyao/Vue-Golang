package controller

import (
	"BTaskServer/model"
	"BTaskServer/util/BLTaskFunc"
	"BTaskServer/util/Tools"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

//易客
//-1:待付款,
//1:已付款,
//2:处理中,
//3:异常,
//4.已完成,
//5:退单中,
//6:已退单,
//7:已退款,
//8:待处理

var YK1typeList []int = []int{424, 423} //易客1type

// 易客1拉单
func GetHigherOrdersYK1(db *gorm.DB) {
	time.Sleep(time.Second * time.Duration(5))
	timer := viper.GetInt("Timer.DownOrderTimer")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("易客1拉取订单recover%s", Tools.GetDateNowFormat(true)))
			go GetHigherOrdersYK1(db)
		}
	}()

	LDpage := 1

	for {
		fmt.Println(fmt.Sprintf("易客1拉取订单%s,page%d", Tools.GetDateNowFormat(true), LDpage))
		resMap := BLTaskFunc.GetOrderListYK1(LDpage, 50, 2)
		if resMap["ret"].(bool) == false {
			fmt.Println("易客1拉取失败，" + resMap["msg"].(string))
			time.Sleep(time.Second)
			continue
		}

		data := resMap["data"].(map[string]interface{})
		infos := data["data"].([]interface{})

		if len(infos) == 0 {
			fmt.Println(fmt.Sprintf("易客1拉取订单完成len=0 %s", Tools.GetDateNowFormat(true)))
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
			goods_ids := valueObj["goodsSN"].(string)
			goods_id, _ := strconv.Atoi(goods_ids)
			if isCodeIn(YK1typeList, goods_id) == false {
				continue
			}

			parmes := valueObj["params"].([]interface{})
			Fparmes := parmes[0].(map[string]interface{})
			shortId := Fparmes["value"].(string)
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

			orderId := int(valueObj["orderSN"].(float64))
			buyNums := valueObj["number"].(string)
			buyNum, _ := strconv.Atoi(buyNums)

			prices := valueObj["price"].(string)
			price, _ := strconv.ParseFloat(prices, 64)

			amounts := valueObj["amount"].(string)
			amount, _ := strconv.ParseFloat(amounts, 64)

			var getTaskType int
			if goods_id == 424 {
				getTaskType = 11
			} else if goods_id == 423 {
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
				SellingPrice: price,
				BuyNumber:    buyNum,
				Amount:       amount,
				StartNum:     int(valueObj["startNum"].(float64)),
				CurrentNum:   0,
				CollectNum:   0,
				GoodsStatus:  2,
				Remark:       valueObj["orderRemark"].(string),
				GoodsName:    valueObj["goodsName"].(string),
				UStatus:      0,
				GetTaskType:  getTaskType,
				PingtaiName:  "易客1",
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

		fmt.Println(fmt.Sprintf("易客1插入新订单 %d 条 %s", len(newRecords), Tools.GetDateNowFormat(true)))

		for _, record := range newRecords {
			fmt.Println(fmt.Sprintf("易客1 拉单 %d", record.OrderId))
		}

		if len(newRecords) > 0 {
			db.Create(&newRecords)
		}

		time.Sleep(time.Second * time.Duration(timer))
	}

}

// 易客1退单
func GetTDOrdersYK1(db *gorm.DB, status int) {
	time.Sleep(time.Second * time.Duration(10))
	timer := viper.GetInt("Timer.UpdateTDOrderTimer")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("易客1退单recover%s", Tools.GetDateNowFormat(true)))
			go GetTDOrdersYK1(db, status)
		}
	}()

	TDpage := 1

	for {
		tx := db.Begin()
		resMap := BLTaskFunc.GetOrderListYK1(TDpage, 50, status)
		fmt.Println(fmt.Sprintf("易客1退单%d订单%s,page%d", status, Tools.GetDateNowFormat(true), TDpage))
		if resMap["ret"].(bool) == false {
			fmt.Println(fmt.Sprintf("易客1退单%d订单失败,%s", status, resMap["msg"].(string)))
			time.Sleep(time.Second)
			continue
		}
		data := resMap["data"].(map[string]interface{})
		infos := data["data"].([]interface{})

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

			orderid := int(valueObj["orderSN"].(float64))
			orderIdList = append(orderIdList, orderid)
		}

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

		fmt.Println(fmt.Sprintf("易客1退单%d订单%s完成", status, Tools.GetDateNowFormat(true)))
		time.Sleep(time.Second * time.Duration(timer))
	}
}
