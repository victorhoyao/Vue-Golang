package controller

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// 领取type ks粉11 ks赞12

type ITaskListController interface {
}

type TaskListController struct {
	DB          *gorm.DB
	RDS         *redis.Client
	TaskTimeOut int
	BfMaxCount  int
}

func NewTaskListController() ITaskListController {
	db := common.GetDB()
	timeout := viper.GetInt("TaskConfig.TaskTimeOut")
	bfMaxCount := viper.GetInt("TaskConfig.BFMaxCount")

	rds := common.GetRedis()

	if err := db.AutoMigrate(&model.TaskList{}); err != nil {
		panic("taskList表迁移失败")
	}
	fmt.Println("taskList表迁移成功")

	// 11ks粉 12快手赞

	// 易客1
	go GetHigherOrdersYK1(db) // 新单
	go GetTDOrdersYK1(db, 5)  //退
	go GetTDOrdersYK1(db, 6)  //退
	go GetTDOrdersYK1(db, 7)  //退

	//通用
	go DelData(db) // 定时清除数据
	go UpdateDoneList01(db)
	go UpdateTaskNoGet(db, rds) // 清除未提交 -暂时改为更新数量

	return TaskListController{DB: db, RDS: rds, TaskTimeOut: timeout, BfMaxCount: bfMaxCount}
}
