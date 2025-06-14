package controller

import (
	"BTaskServer/model"
	"BTaskServer/service"
	"BTaskServer/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskItemController struct{}

func NewTaskItemController() *TaskItemController {
	return &TaskItemController{}
}

// CreateTaskItem handles the creation of a new task item.
func (tc *TaskItemController) CreateTaskItem(c *gin.Context) {
	var requestData struct {
		TaskID        string  `json:"taskId" binding:"required"`
		TaskName      string  `json:"taskName" binding:"required"`
		Supplier      string  `json:"supplier"`
		UserTypes     string  `json:"userTypes"`
		PriceReal     float64 `json:"priceReal"`
		PriceProtocol float64 `json:"priceProtocol"`
		PriceManual   float64 `json:"priceManual"`
		SelectionMode string  `json:"selectionMode"`
		Status        string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		response.Fail(c, nil, "Invalid request data")
		return
	}

	// Set default status if not provided
	if requestData.Status == "" {
		requestData.Status = "active"
	}

	taskItem := model.TaskItem{
		TaskID:        requestData.TaskID,
		TaskName:      requestData.TaskName,
		Supplier:      requestData.Supplier,
		UserTypes:     requestData.UserTypes,
		PriceReal:     requestData.PriceReal,
		PriceProtocol: requestData.PriceProtocol,
		PriceManual:   requestData.PriceManual,
		SelectionMode: requestData.SelectionMode,
		Status:        requestData.Status,
	}

	createdTaskItem, err := service.CreateTaskItem(&taskItem)
	if err != nil {
		response.Fail(c, nil, "Failed to create task item: "+err.Error())
		return
	}

	response.Success(c, gin.H{"taskItem": createdTaskItem}, "Task item created successfully")
}

// GetTaskItems handles retrieving a list of task items.
func (tc *TaskItemController) GetTaskItems(c *gin.Context) {
	taskItems, err := service.GetTaskItems()
	if err != nil {
		response.Fail(c, nil, "Failed to retrieve task items: "+err.Error())
		return
	}
	response.Success(c, gin.H{"taskItems": taskItems}, "Task items retrieved successfully")
}

// UpdateTaskItem handles updating an existing task item.
func (tc *TaskItemController) UpdateTaskItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid task item ID")
		return
	}

	var requestData struct {
		TaskID        string  `json:"taskId"`
		TaskName      string  `json:"taskName"`
		Supplier      string  `json:"supplier"`
		UserTypes     string  `json:"userTypes"`
		PriceReal     float64 `json:"priceReal"`
		PriceProtocol float64 `json:"priceProtocol"`
		PriceManual   float64 `json:"priceManual"`
		SelectionMode string  `json:"selectionMode"`
		Status        string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		response.Fail(c, nil, "Invalid request data")
		return
	}

	taskItem := model.TaskItem{
		TaskID:        requestData.TaskID,
		TaskName:      requestData.TaskName,
		Supplier:      requestData.Supplier,
		UserTypes:     requestData.UserTypes,
		PriceReal:     requestData.PriceReal,
		PriceProtocol: requestData.PriceProtocol,
		PriceManual:   requestData.PriceManual,
		SelectionMode: requestData.SelectionMode,
		Status:        requestData.Status,
	}

	updatedTaskItem, err := service.UpdateTaskItem(uint(id), &taskItem)
	if err != nil {
		response.Fail(c, nil, "Failed to update task item: "+err.Error())
		return
	}

	response.Success(c, gin.H{"taskItem": updatedTaskItem}, "Task item updated successfully")
}

// DeleteTaskItem handles deleting a task item.
func (tc *TaskItemController) DeleteTaskItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid task item ID")
		return
	}

	if err := service.DeleteTaskItem(uint(id)); err != nil {
		response.Fail(c, nil, "Failed to delete task item: "+err.Error())
		return
	}

	response.Success(c, nil, "Task item deleted successfully")
}
