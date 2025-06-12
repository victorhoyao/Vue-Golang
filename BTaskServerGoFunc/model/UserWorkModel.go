package model

type UserWork struct {
	ID         uint   `gorm:"uniqueIndex;primarykey" json:"id" form:"id"`
	CreaTime   string `gorm:"column:creaTime;type:varchar(20);comment:创建时间" json:"creaTime" form:"creaTime"`
	UpdateTime string `gorm:"column:updateTime;type:varchar(20);comment:更新时间" json:"updateTime" form:"updateTime"`
	CountDate  string `gorm:"index;column:countDate;type:varchar(20);comment:统计的日期" json:"countDate" form:"countDate"`
	UserId     uint   `gorm:"index;column:userId;comment:用户id" json:"userId" form:"userId"`
	UserKey    string `gorm:"index;column:userKey;type:varchar(255);not null;comment:用户key" json:"userKey" form:"userKey"`
	UserName   string `gorm:"column:userName;type:varchar(20);not null;comment:账号" json:"userName" form:"userName"`
	// 总
	BliDiggCount      int `gorm:"column:bliDiggCount;not null;default:0;comment:哔哩赞总数量" json:"bliDiggCount" form:"bliDiggCount"`
	BliSLCount        int `gorm:"column:bliSLCount;not null;default:0;comment:哔哩三连总数量" json:"bliSLCount" form:"bliSLCount"`
	BlTBCount         int `gorm:"column:blTBCount;not null;default:0;comment:哔哩投币总数量" json:"blTBCount" form:"blTBCount"`
	BlfenCount        int `gorm:"column:blfenCount;not null;default:0;comment:哔哩粉总数量" json:"blfenCount" form:"blfenCount"`
	BlhuiyuanGouCount int `gorm:"column:blhuiyuanGouCount;not null;default:0;comment:哔哩会员购总数量" json:"blhuiyuanGouCount" form:"blhuiyuanGouCount"`
	BlbofangCount     int `gorm:"column:blbofangCount;not null;default:0;comment:bili播放总数量" json:"blbofangCount" form:"blbofangCount"`
	BlgsfxCount       int `gorm:"column:blgsfxCount;not null;default:0;comment:哔哩高速分享总数量" json:"blgsfxCount" form:"blgsfxCount"`
	BlgsscCount       int `gorm:"column:blgsscCount;not null;default:0;comment:哔哩高速收藏总数量" json:"blgsscCount" form:"blgsscCount"`
	KSDiggCount       int `gorm:"column:KSDiggCount;not null;default:0;comment:KS点赞总数量" json:"KSDiggCount" form:"KSDiggCount"`
	KSSCCount         int `gorm:"column:KSSCCount;not null;default:0;comment:KS收藏总数量" json:"KSSCCount" form:"KSSCCount"`
	KSfenCount        int `gorm:"column:KSfenCount;not null;default:0;comment:KS粉总数量" json:"KSfenCount" form:"KSfenCount"`
	// 成功
	BliDiggSucessCount      int `gorm:"column:bliDiggSucessCount;not null;default:0;comment:哔哩赞成功数量" json:"bliDiggSucessCount" form:"bliDiggSucessCount"`
	BliSLSucessCount        int `gorm:"column:bliSLSucessCount;not null;default:0;comment:哔哩三连成功数量" json:"bliSLSucessCount" form:"bliSLSucessCount"`
	BlTBSucessCount         int `gorm:"column:blTBSucessCount;not null;default:0;comment:哔哩投币成功数量" json:"blTBSucessCount" form:"blTBSucessCount"`
	BlfenSucessCount        int `gorm:"column:blfenSucessCount;not null;default:0;comment:哔哩粉成功数量" json:"blfenSucessCount" form:"blfenSucessCount"`
	BlhuiyuanGouSucessCount int `gorm:"column:blhuiyuanGouSucessCount;not null;default:0;comment:哔哩会员购成功数量" json:"blhuiyuanGouSucessCount" form:"blhuiyuanGouSucessCount"`
	BlbofangSucessCount     int `gorm:"column:blbofangSucessCount;not null;default:0;comment:bili播放成功数量" json:"blbofangSucessCount" form:"blbofangSucessCount"`
	BlgsfxSucessCount       int `gorm:"column:blgsfxSucessCount;not null;default:0;comment:哔哩高速分享成功数量" json:"blgsfxSucessCount" form:"blgsfxSucessCount"`
	BlgsscSucessCount       int `gorm:"column:blgsscSucessCount;not null;default:0;comment:哔哩高速收藏成功数量" json:"blgsscSucessCount" form:"blgsscSucessCount"`
	KSDiggSucessCount       int `gorm:"column:KSDiggSucessCount;not null;default:0;comment:KS点赞成功数量" json:"KSDiggSucessCount" form:"KSDiggSucessCount"`
	KSSCSucessCount         int `gorm:"column:KSSCSucessCount;not null;default:0;comment:KS收藏成功数量" json:"KSSCSucessCount" form:"KSSCSucessCount"`
	KSfenSucessCount        int `gorm:"column:KSfenSucessCount;not null;default:0;comment:KS粉成功数量" json:"KSfenSucessCount" form:"KSfenSucessCount"`
	// 失败
	BliDiggFailCount      int `gorm:"column:bliDiggFailCount;not null;default:0;comment:哔哩赞失败数量" json:"bliDiggFailCount" form:"bliDiggFailCount"`
	BliSLFailCount        int `gorm:"column:bliSLFailCount;not null;default:0;comment:哔哩三连失败数量" json:"bliSLFailCount" form:"bliSLFailCount"`
	BlTBFailCount         int `gorm:"column:blTBFailCount;not null;default:0;comment:哔哩投币失败数量" json:"blTBFailCount" form:"blTBFailCount"`
	BlfenFailCount        int `gorm:"column:blfenFailCount;not null;default:0;comment:哔哩粉失败数量" json:"blfenFailCount" form:"blfenFailCount"`
	BlhuiyuanGouFailCount int `gorm:"column:blhuiyuanGouFailCount;not null;default:0;comment:哔哩会员购失败数量" json:"blhuiyuanGouFailCount" form:"blhuiyuanGouFailCount"`
	BlbofangFailCount     int `gorm:"column:blbofangFailCount;not null;default:0;comment:bili播放失败数量" json:"blbofangFailCount" form:"blbofangFailCount"`
	BlgsfxFailCount       int `gorm:"column:blgsfxFailCount;not null;default:0;comment:哔哩高速分享失败数量" json:"blgsfxFailCount" form:"blgsfxFailCount"`
	BlgsscFailCount       int `gorm:"column:blgsscFailCount;not null;default:0;comment:哔哩高速收藏失败数量" json:"blgsscFailCount" form:"blgsscFailCount"`
	KSDiggFailCount       int `gorm:"column:KSDiggFailCount;not null;default:0;comment:KS点赞失败数量" json:"KSDiggFailCount" form:"KSDiggFailCount"`
	KSSCFailCount         int `gorm:"column:KSSCFailCount;not null;default:0;comment:KS收藏失败数量" json:"KSSCFailCount" form:"KSSCFailCount"`
	KSfenFailCount        int `gorm:"column:KSfenFailCount;not null;default:0;comment:KS粉失败数量" json:"KSfenFailCount" form:"KSfenFailCount"`
}

// TableName 指定表名
func (UserWork) TableName() string {
	return "userWork"
}

// 获取入参
type GetUserWorkQuery struct {
	UserId    int    `json:"userId" form:"userId"`
	CountDate string `json:"countDate" form:"countDate"`
	PageNum   *int   `json:"pageNum" form:"pageNum" binding:"required"`
	PageSize  *int   `json:"pageSize" form:"pageSize" binding:"required"`
}

// 获取自己的入参
type GetMyWorkQuery struct {
	CountDate string `json:"countDate" form:"countDate"`
	PageNum   *int   `json:"pageNum" form:"pageNum" binding:"required"`
	PageSize  *int   `json:"pageSize" form:"pageSize" binding:"required"`
}
