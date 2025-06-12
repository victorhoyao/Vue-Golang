package controller

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"fmt"
	"gorm.io/gorm"
)

type IUserWorkController interface {
}

type UserWorkController struct {
	DB *gorm.DB
}

func NewUserWorkController() IUserWorkController {
	db := common.GetDB()
	if err := db.AutoMigrate(&model.UserWork{}); err != nil {
		panic("userWork表迁移失败")
	}
	fmt.Println("userWork表迁移成功")

	return UserWorkController{DB: db}
}
