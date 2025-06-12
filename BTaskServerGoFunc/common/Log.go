package common

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func InitLog() bool {
	// 输出到日志文件
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)
	if err == nil {
		log.SetOutput(file)         // 输出日志到文件
		log.SetLevel(log.InfoLevel) // 设置日志基本为info最高
		log.SetReportCaller(false)  // 显示行号
		fmt.Println("Log初始化成功")
		return true
	} else {
		fmt.Println("Log初始化失败")
		return false
	}

	// 下面输出，日志级别由低到高，输出情况由上面的日志级别控制
	// 例如：设置日志级别为 ErrorLevel，则 infor、debug、warn不再输出
	//log.Info("info")
	//log.Debug("debug")
	//log.Warn("warn")
	//log.Error("error")
	//log.Panic("panic")
	//log.Fatal("fatal")
}
