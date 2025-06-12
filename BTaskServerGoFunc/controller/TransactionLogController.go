package controller

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"fmt"
	"gorm.io/gorm"
)

type ITransactionLogController interface {
}

type TransactionLogController struct {
	DB *gorm.DB
}

func NewTransactionLogController() ITransactionLogController {
	db := common.GetDB()
	if err := db.AutoMigrate(&model.TransactionLog{}); err != nil {
		panic("transactionLog表迁移失败")
	}

	fmt.Println("transactionLog表迁移成功")

	return TransactionLogController{DB: db}
}
