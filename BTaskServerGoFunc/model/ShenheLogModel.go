package model

type ShengheLog struct {
	ID             uint    `gorm:"uniqueIndex ;primarykey" json:"id" form:"id"`
	UserId         uint    `gorm:"index;column:userId;comment:用户id" json:"userId" from:"userId"`
	UserKey        string  `gorm:"index;column:userKey;type:varchar(255);not null;comment:用户key" json:"userKey" form:"userKey"`
	AllTotal       int     `gorm:"column:allTotal;comment:本轮总审核条数" json:"allTotal" form:"allTotal"`
	EffectiveTotal int     `gorm:"column:effectiveTotal;comment:有效的条数(不偷的)" json:"effectiveTotal" form:"effectiveTotal"`
	DiggTotal      int     `gorm:"column:diggTotal;comment:赞数" json:"diggTotal" form:"diggTotal"`
	FenTotal       int     `gorm:"column:fenTotal;comment:粉数" json:"fenTotal" form:"fenTotal"`
	AddPrice       float64 `gorm:"column:addPrice;type:decimal(11,6);comment:本轮增加金额" json:"addPrice" form:"addPrice"`
	ExamineTime    string  `gorm:"column:examineTime;type:varchar(20);comment:审核时间" json:"examineTime" form:"examineTime"`
}

// TableName 指定表名
func (ShengheLog) TableName() string {
	return "shengheLog"
}
