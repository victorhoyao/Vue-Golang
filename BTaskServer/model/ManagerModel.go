package model

type Manager struct {
	ID         uint   `gorm:"uniqueIndex;primarykey" json:"id" form:"id"`
	ManagerId  uint   `gorm:"index;column:managerId;comment:管理员id" json:"managerId" form:"managerId"`
	UpdateTime string `gorm:"column:updateTime;type:varchar(20);comment:更新时间" json:"updateTime" form:"updateTime"`
	//BliDiggPrice float64 `gorm:"column:bliDiggPrice;type:decimal(11,6);comment:哔哩赞单价" json:"bliDiggPrice" form:"bliDiggPrice"`
	//BliSLPrice        float64 `gorm:"column:bliSLPrice;type:decimal(11,6);comment:哔哩三连单价" json:"bliSLPrice" form:"bliSLPrice"`
	//BlTBPrice         float64 `gorm:"column:blTBPrice;type:decimal(11,6);comment:哔哩投币单价" json:"blTBPrice" form:"blTBPrice"`
	//BlfenPrice float64 `gorm:"column:blfenPrice;type:decimal(11,6);comment:哔哩粉单价" json:"blfenPrice" form:"blfenPrice"`
	//BlhuiyuanGouPrice float64 `gorm:"column:blhuiyuanGouPrice;type:decimal(11,6);comment:哔哩会员购单价" json:"blhuiyuanGouPrice" form:"blhuiyuanGouPrice"`
	//BlbofangPrice     float64 `gorm:"column:blbofangPrice;type:decimal(11,6);comment:bili播放单价" json:"blbofangPrice" form:"blbofangPrice"`
	//BlgsfxPrice       float64 `gorm:"column:blgsfxPrice;type:decimal(11,6);comment:哔哩高速分享单价" json:"blgsfxPrice" form:"blgsfxPrice"`
	//BlgsscPrice       float64 `gorm:"column:blgsscPrice;type:decimal(11,6);comment:哔哩高速收藏单价" json:"blgsscPrice" form:"blgsscPrice"`
	KSDiggPrice float64 `gorm:"column:KSDiggPrice;type:decimal(11,6);comment:KS点赞单价" json:"KSDiggPrice" form:"KSDiggPrice"`
	//KSSCPrice         float64 `gorm:"column:KSSCPrice;type:decimal(11,6);comment:KS收藏单价" json:"KSSCPrice" form:"KSSCPrice"`
	KSfenPrice float64 `gorm:"column:KSfenPrice;type:decimal(11,6);comment:KS粉单价" json:"KSfenPrice" form:"KSfenPrice"`
	//DyBfPrice float64 `gorm:"column:DyBfPrice;type:decimal(11,6);not null;default:0;comment:DY播放单价" json:"DyBfPrice" form:"DyBfPrice"`
	//DyBfhyPrice float64 `gorm:"column:DyBfhyPrice;type:decimal(11,6);comment:DY播放dy单价620" json:"DyBfhyPrice" form:"DyBfhyPrice"`
	ShTcGl int `gorm:"column:shTcGl;comment:假审核概率,百分之多少" json:"shTcGl" form:"shTcGl"`
}

// TableName 指定表名
func (Manager) TableName() string {
	return "manager"
}

// 更新单价入参
type SetPriceJson struct {
	Type  string  `json:"type" form:"type" binding:"required,oneof=快手赞 快手粉"`
	Price float64 `json:"price" json:"price" binding:"required"`
}

// 更新偷吃审核入参
type SetTcGlJson struct {
	ShTcGl *int `json:"shTcGl" form:"shTcGl" binding:"required"`
}
