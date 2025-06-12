package controller

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"fmt"
	"gorm.io/gorm"
)

type ITaskLogController interface {
}

type TaskLogController struct {
	DB *gorm.DB
}

func NewTaskLogController() ITaskLogController {
	db := common.GetDB()
	if err := db.AutoMigrate(&model.TaskLog{}); err != nil {
		panic("taskLog表迁移失败")
	}
	fmt.Println("taskLog表迁移成功")

	if err1 := db.AutoMigrate(&model.ShengheLog{}); err1 != nil {
		panic("shengheLog表迁移失败")
	}
	fmt.Println("shengheLog表迁移成功")

	go Shenghe(db)

	return TaskLogController{DB: db}
}
