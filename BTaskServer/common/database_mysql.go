package common

import (
	"BTaskServer/global"
	"BTaskServer/model"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql" // gorm mysql 驱动包
	"gorm.io/gorm"         // gorm3
	"gorm.io/gorm/schema"
)

func InitDB() bool {
	// MySQL 配置信息
	username := viper.GetString("datasouce_mysql.admin")    // 账号
	password := viper.GetString("datasouce_mysql.password") // 密码
	host := viper.GetString("datasouce_mysql.host")         // 地址

	// 端口
	port, _ := strconv.Atoi(viper.GetString("datasouce_mysql.port"))

	DBname := viper.GetString("datasouce_mysql.database") // 数据库名称

	timeout := "10s" // 连接超时，10秒

	//dsn1 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, DBname)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, DBname, timeout)

	// Open 连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用后，`User` 表将是 `user`
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	})
	if err != nil {
		fmt.Println("数据库连接失败    ", err)
		return false
	}
	fmt.Println("数据库连接成功    ")

	err = db.AutoMigrate(&model.User{}, &model.TransactionLog{}, &model.Manager{}, &model.TaskList{}, &model.TaskLog{}, &model.Supplier{}, &model.TaskItem{}, &model.TaskDistribution{})
	if err != nil {
		fmt.Println("表迁移失败 ", err)
		return false
	}
	fmt.Println("user表迁移成功")
	fmt.Println("transactionLog表迁移成功")
	fmt.Println("manager表迁移成功 ")
	fmt.Println("taskList表迁移成功")
	fmt.Println("taskLog表迁移成功")
	fmt.Println("supplier表迁移成功")
	fmt.Println("taskItem表迁移成功")
	fmt.Println("taskDistribution表迁移成功")

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.GVA_DB = db // Assign to global.GVA_DB instead of common.DB
	return true
}

// GetDB function is no longer needed if global.GVA_DB is used directly
/*
func GetDB() *gorm.DB {
	return DB
}
*/
