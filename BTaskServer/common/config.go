package common

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitConfig() bool {
	workDir, _ := os.Getwd()
	fmt.Println(workDir)
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("配置文件读取失败")
		return false
	}
	fmt.Println("配置文件读取成功")
	return true
}
