package controller

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"fmt"
	"gorm.io/gorm"
)

type IUserController interface {
}

type UserController struct {
	DB *gorm.DB
}

func NewUserController() IUserController {
	userDB := common.GetDB()
	if err := userDB.AutoMigrate(&model.User{}); err != nil {
		panic("user表迁移失败")
	}
	fmt.Println("user表迁移成功")

	return UserController{DB: userDB}
}
