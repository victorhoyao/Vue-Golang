package common

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql" // gorm mysql 驱动包
	"gorm.io/gorm"         // gorm3
	"strconv"
	"time"
)

var DB *gorm.DB

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败.")
		return false
	}
	fmt.Println("数据库连接成功")

	DB = db

	s, err := DB.DB()
	s.SetMaxOpenConns(100)
	s.SetMaxIdleConns(10)
	s.SetConnMaxLifetime(10 * time.Minute)
	s.SetConnMaxIdleTime(5 * time.Minute)
	return true
}

func GetDB() *gorm.DB {
	return DB
}
