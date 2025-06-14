package model

type User struct {
	ID          uint    `gorm:"uniqueIndex;primarykey" json:"id" form:"id"`
	UserName    string  `gorm:"column:userName;type:varchar(20);not null;comment:账号" json:"userName" form:"userName"`
	PassWord    string  `gorm:"column:passWord;type:varchar(255);not null;comment:密码" json:"-" form:"-"`
	Status      int     `gorm:"column:status;comment:用户状态0禁用1启用" json:"status" form:"status"`
	Authority   int     `gorm:"column:authority;comment:权限编码0普通用户1管理员" json:"authority" form:"authority"`
	AddUserId   uint    `gorm:"column:addUserId;comment:添加用户id" json:"addUserId" form:"addUserId"`
	UserKey     string  `gorm:"index;column:userKey;type:varchar(255);not null;comment:用户key" json:"userKey" form:"userKey"`
	Money       float64 `gorm:"column:money;type:decimal(11,6);comment:金额" json:"money" form:"money"`
	TranType    string  `gorm:"column:tranType;type:varchar(20);comment:提现账号类型" json:"tranType" form:"tranType"`
	TranAccount string  `gorm:"column:tranAccount;type:varchar(255);comment:提现账号" json:"tranAccount" form:"tranAccount"`
	TranName    string  `gorm:"column:tranName;type:varchar(50);comment:提现账号姓名" json:"tranName" form:"tranName"`
	CreateTime  string  `gorm:"column:createTime;type:varchar(20);comment:创建时间" json:"createTime" form:"createTime"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// 用户登录入参
type UserLogin struct {
	UserName string `json:"userName" form:"userName" binding:"required,max=20"`
	PassWord string `json:"passWord" form:"passWord" binding:"required,max=20"`
}

// 添加用户入参
type AddUser struct {
	UserName string `json:"userName" form:"userName" binding:"required,max=20"`
	PassWord string `json:"passWord" form:"passWord" binding:"required,max=20"`
}

// 添加管理员入参
type AddManager struct {
	UserName string `json:"userName" form:"userName" binding:"required,max=20"`
	PassWord string `json:"passWord" form:"passWord" binding:"required,max=20"`
}

// 修改用户密码入参
type EditUserPass struct {
	UserId      uint   `json:"userId" form:"userId" binding:"required"`
	NewPassword string `json:"newPassword" form:"newPassword" binding:"required"`
}

// 删除用户入参
type DelUser struct {
	UserId uint `json:"userId" form:"userId" binding:"required"`
}

// 设置提现账号入参
type SetTranJson struct {
	AccountType string `json:"accountType" form:"accountType" binding:"required,oneof=支付宝 U提现 银行卡"`
	Account     string `json:"account" form:"account" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
}
