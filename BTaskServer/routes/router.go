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
	r.POST("register", UserController.Register)
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
	User.POST("changePassword", UserController.ChangePassword) // 修改用户密码

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

	// Supplier Management routes
	SupplierController := controller.NewSupplierController()
	Supplier := r.Group("Supplier")
	Supplier.Use(middleware.AuthMiddleware())
	Supplier.POST("add", SupplierController.AddSupplier)
	Supplier.GET("list", SupplierController.GetSuppliers)
	Supplier.PUT("update/:id", SupplierController.UpdateSupplier)
	Supplier.DELETE("delete/:id", SupplierController.DeleteSupplier)

	// Task Item Management routes
	TaskItemController := controller.NewTaskItemController()
	TaskItem := r.Group("TaskItem")
	TaskItem.Use(middleware.AuthMiddleware())
	TaskItem.POST("add", TaskItemController.CreateTaskItem)
	TaskItem.GET("list", TaskItemController.GetTaskItems)
	TaskItem.PUT("update/:id", TaskItemController.UpdateTaskItem)
	TaskItem.DELETE("delete/:id", TaskItemController.DeleteTaskItem)

	// Task Distribution Management routes
	TaskDistributionController := controller.NewTaskDistributionController()
	TaskDistribution := r.Group("TaskDistribution")
	TaskDistribution.Use(middleware.AuthMiddleware())
	TaskDistribution.POST("create", TaskDistributionController.CreateTaskDistribution)
	TaskDistribution.GET("list", TaskDistributionController.GetTaskDistributions)
	TaskDistribution.GET("get/:id", TaskDistributionController.GetTaskDistribution)
	TaskDistribution.PUT("update/:id", TaskDistributionController.UpdateTaskDistribution)
	TaskDistribution.DELETE("delete/:id", TaskDistributionController.DeleteTaskDistribution)
	TaskDistribution.POST("activate/:id", TaskDistributionController.ActivateTaskDistribution)
	TaskDistribution.GET("by-task-item/:taskItemId", TaskDistributionController.GetTaskDistributionsByTaskItem)
	TaskDistribution.GET("summary", TaskDistributionController.GetDistributionSummary)

	return r
}
