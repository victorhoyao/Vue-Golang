package model

type TaskList struct {
	ID           uint    `gorm:"uniqueIndex ;primarykey" json:"id" form:"id"`
	DownTime     string  `gorm:"column:downTime;type:varchar(20);comment:订单拉取时间" json:"downTime" form:"downTime"`
	LastGetTime  string  `gorm:"column:lastGetTime;varchar(20);comment:最后一次领取时间" json:"lastGetTime" form:"lastGetTime"`
	DoneTime     string  `gorm:"column:doneTime;varchar(20);comment:订单完成时间" json:"doneTime" form:"doneTime"`
	VideoLink    string  `gorm:"column:videoLink;comment:视频长链接" json:"videoLink" form:"videoLink"`
	ShortLink    string  `gorm:"column:shortLink;comment:视频短链接" json:"shortLink" form:"shortLink"`
	VideoId      string  `gorm:"column:videoId;type:varchar(255);comment:视频id" json:"videoId" form:"videoId"`
	Uid          string  `gorm:"column:uid;type:varchar(255);comment:uid" json:"uid" form:"uid"`
	OrderId      int     `gorm:"index;column:orderId;comment:上级_订单id" json:"orderId" form:"orderId"`
	GoodsId      int     `gorm:"index;column:goodsId;comment:上级_goods_id" json:"goodsId" form:"goodsId"`
	SellingPrice float64 `gorm:"column:sellingPrice;type:decimal(11,6);comment:上级_商品单价" json:"sellingPrice" form:"sellingPrice"`
	BuyNumber    int     `gorm:"index;column:buyNumber;comment:上级_购买数量" json:"buyNumber" form:"buyNumber"`
	Amount       float64 `gorm:"column:amount;type:decimal(11,6);comment:上级_商品总价" json:"amount" form:"amount"`
	StartNum     int     `gorm:"column:startNum;comment:上级_开始进度" json:"startNum" form:"startNum"`
	CurrentNum   int     `gorm:"column:currentNum;comment:上级_当前进度" json:"currentNum" form:"currentNum"`
	CollectNum   int     `gorm:"index;column:collectNum;comment:领取数量" json:"collectNum" form:"collectNum"`
	GoodsStatus  int     `gorm:"column:goodsStatus;comment:上级_商品状态 已付款|处理中.." json:"goodsStatus" form:"goodsStatus"`
	Remark       string  `gorm:"column:remark;comment:上级_订单备注" json:"remark" form:"remark"`
	GoodsName    string  `gorm:"column:goodsName;type:varchar(20);comment:上级_商品名称" json:"goodsName" form:"goodsName"`
	UStatus      int     `gorm:"index;column:UStatus;comment:订单状态 0处理中 1已完成 2退单 3播放任务完成并重置" json:"UStatus" form:"UStatus"`
	GetTaskType  int     `gorm:"index;column:getTaskType;comment:领取的type" json:"getTaskType" form:"getTaskType"`
	PingtaiName  string  `gorm:"column:pingtaiName;type:varchar(20);comment:任务来源" json:"pingtaiName" form:"pingtaiName"`
	ErrCount     int     `gorm:"column:errCount;default:0;comment:查询数量错误次数" json:"errCount" form:"errCount"`
}

// TableName 指定表名
func (TaskList) TableName() string {
	return "taskList"
}

// 领取入参
type GetTaskQuery struct {
	Time      int    `json:"time" form:"time" binding:"required"`
	Jcc       string `json:"jcc" form:"jcc" binding:"required"`
	UserKey   string `json:"userKey" form:"userKey" binding:"required"`
	BlAccount string `json:"blAccount" form:"blAccount" binding:"required"`
	Type      int    `json:"type" form:"type" binding:"required"`
}

// 领取任务出参
type GetTaskModel struct {
	ID        int    `json:"id" form:"id"`
	VideoLink string `json:"videoLink" form:"videoLink"`
	ShortLink string `json:"shortLink" form:"shortLink"`
	VideoId   string `json:"videoId" form:"videoId"`
	Uid       string `json:"uid" form:"uid"`
	GoodsName string `json:"goodsName" form:"goodsName"`
}

// 提交任务入参
type SubmitTaskQuery struct {
	ID int `json:"id" form:"id" binding:"required"`
}

// 任务列表搜索
type FindTaskQuery struct {
	PageNum  *int   `json:"pageNum" form:"pageNum" binding:"required"`
	PageSize *int   `json:"pageSize" form:"pageSize" binding:"required"`
	OrderId  string `json:"orderId" form:"orderId" binding:"required"`
}
