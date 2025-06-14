package controller

import (
	"BTaskServer/global"
	"BTaskServer/model"
	"BTaskServer/util/Tools"
	"BTaskServer/util/response"
	"BTaskServer/util/validatorTool"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IManagerController interface {
	SetPrice(c *gin.Context)
	GetPrice(c *gin.Context)
	SetTcGl(c *gin.Context)
}

type ManagerController struct {
	DB *gorm.DB
}

var managerList []model.Manager

func (m ManagerController) SetTcGl(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)
	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	var setTcGlJson model.SetTcGlJson
	if !validatorTool.ValidatorJson[*model.SetTcGlJson](c, &setTcGlJson) {
		return
	}

	var manager model.Manager
	res := m.DB.Where("managerId = ?", usermodel.ID).Limit(1).Find(&manager)
	if res.Error != nil {
		response.ServerBad(c, nil, "查询数据失败")
		return
	}

	if res.RowsAffected == 0 {
		newManager := model.Manager{
			ManagerId:  usermodel.ID,
			UpdateTime: Tools.GetDateNowFormat(true),
			ShTcGl:     *setTcGlJson.ShTcGl,
		}

		if res1 := m.DB.Create(&newManager); res1.Error != nil {
			response.ServerBad(c, nil, "修改失败")
			return
		}
	} else {
		res2 := m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			"shTcGl":     setTcGlJson.ShTcGl,
			"updateTime": Tools.GetDateNowFormat(true),
		})
		if res2.Error != nil {
			response.ServerBad(c, nil, "修改失败")
			return
		}
	}

	cacheManagerList()
	response.Success(c, nil, "修改成功")
}

func (m ManagerController) GetPrice(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	var manager model.Manager
	if usermodel.Authority == 1 {
		res := m.DB.Where("managerId = ?", usermodel.ID).Limit(1).Find(&manager)
		if res.Error != nil {
			response.ServerBad(c, nil, "查询失败")
			return
		}
		if res.RowsAffected == 0 {
			response.ServerBad(c, nil, "查询失败")
			return
		}
	} else {
		res1 := m.DB.Where("managerId = ?", usermodel.AddUserId).Limit(1).Find(&manager)
		if res1.Error != nil {
			response.ServerBad(c, nil, "查询失败")
			return
		}
		if res1.RowsAffected == 0 {
			response.ServerBad(c, nil, "查询失败")
			return
		}
	}

	response.Success(c, gin.H{"result": manager}, "获取成功")
}

func (m ManagerController) SetPrice(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	var setPriceJson model.SetPriceJson
	if !validatorTool.ValidatorJson[*model.SetPriceJson](c, &setPriceJson) {
		return
	}

	var manager model.Manager
	res := m.DB.Where("managerId = ?", usermodel.ID).Limit(1).Find(&manager)
	if res.Error != nil {
		response.ServerBad(c, nil, "查询单价数据失败")
		return
	}

	if res.RowsAffected == 0 {
		var newManager model.Manager
		switch setPriceJson.Type {

		case "快手粉":
			newManager = model.Manager{
				ManagerId:  usermodel.ID,
				UpdateTime: Tools.GetDateNowFormat(true),
				KSfenPrice: setPriceJson.Price,
			}
		case "快手赞":
			newManager = model.Manager{
				ManagerId:   usermodel.ID,
				UpdateTime:  Tools.GetDateNowFormat(true),
				KSDiggPrice: setPriceJson.Price,
			}
		}

		if res1 := m.DB.Create(&newManager); res1.Error != nil {
			response.ServerBad(c, nil, "修改失败")
			return
		}
	} else {
		var res2 *gorm.DB
		switch setPriceJson.Type {
		case "快手赞":
			res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
				"KSDiggPrice": setPriceJson.Price,
				"updateTime":  Tools.GetDateNowFormat(true),
			})
		//case "哔哩三连":
		//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
		//		"bliSLPrice": setPriceJson.Price,
		//		"updateTime": Tools.GetDateNowFormat(true),
		//	})
		case "快手粉":
			res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
				"KSfenPrice": setPriceJson.Price,
				"updateTime": Tools.GetDateNowFormat(true),
			})
			//case "哔哩投币":
			//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			//		"blTBPrice":  setPriceJson.Price,
			//		"updateTime": Tools.GetDateNowFormat(true),
			//	})
			//case "哔哩会员购":
			//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			//		"blhuiyuanGouPrice": setPriceJson.Price,
			//		"updateTime":        Tools.GetDateNowFormat(true),
			//	})
			//case "bili播放":
			//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			//		"blbofangPrice": setPriceJson.Price,
			//		"updateTime":    Tools.GetDateNowFormat(true),
			//	})
			//case "哔哩高速收藏":
			//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			//		"blgsscPrice": setPriceJson.Price,
			//		"updateTime":  Tools.GetDateNowFormat(true),
			//	})
			//case "哔哩高速分享":
			//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			//		"blgsfxPrice": setPriceJson.Price,
			//		"updateTime":  Tools.GetDateNowFormat(true),
			//	})
			//case "KS收藏":
			//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			//		"KSSCPrice":  setPriceJson.Price,
			//		"updateTime": Tools.GetDateNowFormat(true),
			//	})
			//case "KS赞":
			//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			//		"KSDiggPrice": setPriceJson.Price,
			//		"updateTime":  Tools.GetDateNowFormat(true),
			//	})
			//case "KS粉":
			//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			//		"KSfenPrice": setPriceJson.Price,
			//		"updateTime": Tools.GetDateNowFormat(true),
			//	})
			//case "DY播放":
			//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			//		"DyBfPrice":  setPriceJson.Price,
			//		"updateTime": Tools.GetDateNowFormat(true),
			//	})
			//case "DY播放hy":
			//
			//	res2 = m.DB.Model(&model.Manager{}).Where("managerId = ?", usermodel.ID).Updates(map[string]interface{}{
			//		"DyBfhyPrice": setPriceJson.Price,
			//		"updateTime":  Tools.GetDateNowFormat(true),
			//	})
		}

		if res2.Error != nil {
			response.ServerBad(c, nil, "修改失败")
			return
		}
	}

	cacheManagerList()
	response.Success(c, nil, "修改成功")
}

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
	return resMap
}

func cacheManagerList() bool {
	res := global.GVA_DB.Find(&managerList)
	if res.Error != nil {
		return false
	}
	return true
}

func getCacelManager(managerId uint) (model.Manager, bool) {
	for _, v := range managerList {
		if v.ManagerId == managerId {
			return v, true
		}
	}
	return model.Manager{}, false
}

func getCacelGaiTc() model.Manager {
	var manager model.Manager
	global.GVA_DB.Where("managerId = ?", 1).Limit(1).Find(&manager)
	return manager
}

func NewManagerController() IManagerController {
	return ManagerController{DB: global.GVA_DB}
}
