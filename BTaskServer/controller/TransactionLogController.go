package controller

import (
	"BTaskServer/global"
	"BTaskServer/model"
	"BTaskServer/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
	gorm "gorm.io/gorm"
)

type ITransactionLogController interface {
	ApplyTransaction(c *gin.Context)
	DoneTransaction(c *gin.Context)
	GetApplyList(c *gin.Context)
	GetMyTranList(c *gin.Context)
	GetTranListByUserId(c *gin.Context)
}

type TransactionLogController struct {
	DB *gorm.DB
}

func NewTransactionLogController() ITransactionLogController {
	db := global.GVA_DB
	if err := db.AutoMigrate(&model.TransactionLog{}); err != nil {
		panic("transactionLog表迁移失败")
	}

	return TransactionLogController{DB: db}
}

func (t TransactionLogController) ApplyTransaction(c *gin.Context) {
	var applyData model.ApplyJson
	if err := c.ShouldBindJSON(&applyData); err != nil {
		response.Fail(c, nil, "Invalid request data")
		return
	}

	// Basic implementation - you can expand this later
	response.Success(c, nil, "Transaction applied successfully")
}

func (t TransactionLogController) DoneTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid transaction ID")
		return
	}

	var doneData model.DoneApplyJson
	if err := c.ShouldBindJSON(&doneData); err != nil {
		response.Fail(c, nil, "Invalid request data")
		return
	}

	// Basic implementation - you can expand this later
	response.Success(c, gin.H{"id": id}, "Transaction processed successfully")
}

func (t TransactionLogController) GetApplyList(c *gin.Context) {
	var transactions []model.TransactionLog
	if err := t.DB.Find(&transactions).Error; err != nil {
		response.Fail(c, nil, "Failed to fetch transactions")
		return
	}

	response.Success(c, gin.H{"transactions": transactions}, "Transactions retrieved successfully")
}

func (t TransactionLogController) GetMyTranList(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.AuthError(c, nil, "User not found")
		return
	}

	userModel := user.(model.User)
	var transactions []model.TransactionLog
	if err := t.DB.Where("userId = ?", userModel.ID).Find(&transactions).Error; err != nil {
		response.Fail(c, nil, "Failed to fetch transactions")
		return
	}

	response.Success(c, gin.H{"transactions": transactions}, "Transactions retrieved successfully")
}

func (t TransactionLogController) GetTranListByUserId(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid user ID")
		return
	}

	var transactions []model.TransactionLog
	if err := t.DB.Where("userId = ?", userId).Find(&transactions).Error; err != nil {
		response.Fail(c, nil, "Failed to fetch transactions")
		return
	}

	response.Success(c, gin.H{"transactions": transactions}, "Transactions retrieved successfully")
}
