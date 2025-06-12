package model

type TransactionLog struct {
	ID          uint    `gorm:"uniqueIndex;primarykey" json:"id" form:"id"`
	UserId      uint    `gorm:"column:userId;comment:用户id" json:"userId" from:"userId"`
	AddUserId   uint    `gorm:"index;column:addUserId;comment:添加用户id" json:"addUserId" form:"addUserId"`
	OrderId     string  `gorm:"column:orderId;type:varchar(255);comment:提现订单号" json:"orderId" form:"orderId"`
	CreateTime  string  `gorm:"column:createTime;type:varchar(20);comment:创建时间" json:"createTime" form:"createTime"`
	DoneTime    string  `gorm:"column:doneTime;type:varchar(20);comment:完成时间" json:"doneTime" form:"doneTime"`
	Applyprice  float64 `gorm:"column:applyprice;type:decimal(11,6);comment:申请金额" json:"applyprice" form:"applyprice"`
	OldPrice    float64 `gorm:"column:oldPrice;type:decimal(11,6);comment:提现前余额" json:"oldPrice" form:"oldPrice"`
	NewPrice    float64 `gorm:"column:newPrice;type:decimal(11,6);comment:提现后余额" json:"newPrice" form:"newPrice"`
	TranType    string  `gorm:"column:tranType;type:varchar(20);comment:提现账号类型" json:"tranType" form:"tranType"`
	TranAccount string  `gorm:"column:tranAccount;type:varchar(255);comment:提现账号" json:"tranAccount" form:"tranAccount"`
	TranName    string  `gorm:"column:tranName;type:varchar(50);comment:提现账号姓名" json:"tranName" form:"tranName"`
	Status      int     `gorm:"column:status;comment:状态0审核中1审核成功2审核失败" json:"status" form:"status"`
	Desc        string  `gorm:"column:desc;type:varchar(255);comment:说明" json:"desc" form:"desc"`
}

// TableName 指定表名
func (TransactionLog) TableName() string {
	return "transactionLog"
}

// 提现入参
type ApplyJson struct {
	Applyprice float64 `json:"applyprice" form:"applyprice" binding:"required"`
}

// 审核提现入参
type DoneApplyJson struct {
	Status int    `json:"status" form:"status" binding:"required,oneof=1 2"`
	Desc   string `json:"desc" form:"desc"`
}
