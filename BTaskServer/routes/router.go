package routes

import (
	"BTaskServer/controller"
	"BTaskServer/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())

	UserController := controller.NewUserController()
	r.POST("login", UserController.Login)
	r.POST("addManager", UserController.AddManager)
	ConfigController := controller.NewConfigController()
	r.GET("getApplicationConfig", ConfigController.GetApplicationConfig)

	// 管理员接口和用户接口
	User := r.Group("User")
	User.Use(middleware.AuthMiddleware())
	User.POST("addUser", UserController.AddUser)       // 管理员
	User.POST("editPass", UserController.EditUserPass) // 管理员
	User.GET("userList", UserController.GetUserList)   // 管理员
	User.GET("getMyInfo", UserController.MyInfo)
	User.GET("delUser", UserController.DelUser) // 管理员
	User.POST("setTran", UserController.SetTran)
	User.POST("setTranByUser/:id", UserController.ManagerSetTranByUserId) // 管理员
	User.GET("findUser", UserController.FindUser)

	TransactionController := controller.NewTransactionLogController()
	Tran := r.Group("Tran")
	Tran.Use(middleware.AuthMiddleware())
	Tran.POST("applyTransaction", TransactionController.ApplyTransaction)
	Tran.POST("doneTransaction/:id", TransactionController.DoneTransaction) // 管理员
	Tran.GET("getApplyList", TransactionController.GetApplyList)            // 管理员
	Tran.GET("getMyTranList", TransactionController.GetMyTranList)
	Tran.GET("getTranListByUserId/:id", TransactionController.GetTranListByUserId) // 管理员

	ManagerController := controller.NewManagerController()
	Manager := r.Group("Manager")
	Manager.Use(middleware.AuthMiddleware())
	Manager.POST("setPrice", ManagerController.SetPrice)
	Manager.GET("getPrice", ManagerController.GetPrice)
	Manager.POST("setTcGl", ManagerController.SetTcGl)

	taskListController := controller.NewTaskListController()
	taskList := r.Group("taskList")
	taskList.GET("getTask", taskListController.GetTask)
	taskList.GET("submitTask", taskListController.SubmitTask)
	taskList.GET("getTaskList", middleware.AuthMiddleware(), taskListController.GetTaskList)
	taskList.GET("getTaskListByKey", middleware.AuthMiddleware(), taskListController.GetTaskListByKey)

	taskLogController := controller.NewTaskLogController()
	taskLog := r.Group("taskLog")
	taskLog.Use(middleware.AuthMiddleware())
	taskLog.GET("getTaskLogList", taskLogController.GetTaskLogList)
	taskLog.GET("getTaskLogListByKey", taskLogController.GetTaskLogListByKey)
	taskLog.GET("getMyTaskLogList", taskLogController.GetMyTaskLogList)
	taskLog.GET("getTaskLogListById/:id", taskLogController.GetTaskLogListById)
	taskLog.GET("getTaskLogCount", taskLogController.GetTaskLogCount)

	return r
}
