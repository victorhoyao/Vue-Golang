package controller

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"fmt"
	"gorm.io/gorm"
)

type IManagerController interface {
}

type ManagerController struct {
	DB *gorm.DB
}

var managerList []model.Manager

// 获取单价
func GetPriceByType(addUserId uint, priceType int) map[string]interface{} {
	resMap := make(map[string]interface{})

	manager, b := getCacelManager(addUserId)
	if b == false {
		resMap["ret"] = false
		resMap["price"] = -1
		resMap["msg"] = "单价查询失败"
		return resMap
	}

	var price float64
	switch priceType {
	case 11:
		price = manager.KSfenPrice
	case 12:
		price = manager.KSDiggPrice
	}

	if price == 0 {
		resMap["ret"] = false
		resMap["price"] = -1
		resMap["msg"] = "该订单类型未设置单价,无法领取,请联系管理员"
		return resMap
	}

	resMap["ret"] = true
	resMap["price"] = price
	resMap["msg"] = ""
	return resMap
}

func cacheManagerList() bool {
	db := common.GetDB()
	managerList = []model.Manager{}
	res := db.Where("1=1").Find(&managerList)
	if res.Error != nil {
		fmt.Println("缓存managerList失败")
		return false
	}
	return true
}

func getCacelManager(managerId uint) (model.Manager, bool) {
	var resmanager model.Manager
	flag := false
	for _, manager := range managerList {
		if manager.ManagerId == managerId {
			resmanager = manager
			flag = true
			break
		}
	}

	return resmanager, flag
}

func getCacelGaiTc() model.Manager {
	return managerList[0]
}

func NewManagerController() IManagerController {
	db := common.GetDB()
	if err := db.AutoMigrate(&model.Manager{}); err != nil {
		panic("manager表迁移失败")
	}
	fmt.Println("manager表迁移成功")

	cacheRet := cacheManagerList()
	if cacheRet == false {
		panic("缓存managerList失败")
	}

	return ManagerController{DB: db}
}
