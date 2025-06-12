package model

type TaskLog struct {
	ID          uint    `gorm:"uniqueIndex;primarykey" json:"id" form:"id"`
	CollectTime string  `gorm:"column:collectTime;type:varchar(20);comment:领取时间" json:"collectTime" form:"collectTime"`
	SumbmitTime string  `gorm:"index;column:sumbmitTime;type:varchar(20);comment:提交时间" json:"sumbmitTime" form:"sumbmitTime"`
	ExamineTime string  `gorm:"column:examineTime;type:varchar(20);comment:审核时间" json:"examineTime" form:"examineTime"`
	OrderNum    string  `gorm:"column:orderNum;type:varchar(255);comment:订单编号" json:"orderNum" form:"orderNum"`
	Account     string  `gorm:"index;column:Account;type:varchar(255);comment:领取账号" json:"Account" form:"blAccount"`
	VideoLink   string  `gorm:"column:videoLink;comment:上级_视频链接" json:"videoLink" form:"videoLink"`
	OrderId     int     `gorm:"index;column:orderId;comment:上级_订单id" json:"orderId" form:"orderId"`
	GoodsId     int     `gorm:"index;column:goodsId;comment:上级_goods_id" json:"goodsId" form:"goodsId"`
	GoodsName   string  `gorm:"column:goodsName;type:varchar(20);comment:上级_商品名称" json:"goodsName" form:"goodsName"`
	UserKey     string  `gorm:"index;column:userKey;type:varchar(255);not null;comment:用户key" json:"userKey" form:"userKey"`
	ManagerId   uint    `gorm:"index;column:managerId;comment:管理员id" json:"managerId" form:"managerId"`
	Price       float64 `gorm:"column:price;type:decimal(11,6);comment:任务单价" json:"price" json:"price"`
	Status      int     `gorm:"index;column:status;comment:状态 0未提交 1已提交等待审核 2过期作废 3审核通过 4审核不通过" json:"status" form:"status"`
	Remark      string  `gorm:"column:remark;type:varchar(200);comment:订单备注" json:"remark" form:"remark"`
	BfCount     int     `gorm:"index;column:bfCount;comment:播放任务领取次数" json:"bfCount" form:"bfCount"`
	GetTaskType int     `gorm:"index;column:getTaskType;comment:领取的type" json:"getTaskType" form:"getTaskType"`
	PingtaiName string  `gorm:"column:pingtaiName;type:varchar(20);comment:任务来源" json:"pingtaiName" form:"pingtaiName"`
}

// TableName 指定表名
func (TaskLog) TableName() string {
	return "taskLog"
}

// 获取任务统计数量入参
type GetTaskLogCountQuery struct {
	Date string `json:"date" form:"date" binding:"required"`
	Type int    `json:"type" form:"type" binding:"required"`
}

// 获取任务统计数量出参
type GetTaskLogCountModel struct {
	WorkCount   int `json:"workCount" form:"workCount"`
	SucessCount int `json:"sucessCount" form:"sucessCount"`
	FailCount   int `json:"failCount" form:"failCount"`
}
