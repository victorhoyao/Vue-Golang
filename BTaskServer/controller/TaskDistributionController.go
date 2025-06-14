package controller

import (
	"BTaskServer/model"
	"BTaskServer/service"
	"BTaskServer/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TaskDistributionController handles task distribution related requests
type TaskDistributionController struct {
	service *service.TaskDistributionService
}

// NewTaskDistributionController creates a new TaskDistributionController
func NewTaskDistributionController() *TaskDistributionController {
	return &TaskDistributionController{
		service: service.NewTaskDistributionService(),
	}
}

// CreateTaskDistribution handles creating a new task distribution
func (tc *TaskDistributionController) CreateTaskDistribution(c *gin.Context) {
	var req model.TaskDistributionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, nil, "Invalid request data: "+err.Error())
		return
	}

	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		response.AuthError(c, nil, "User not found")
		return
	}
	userModel := user.(model.User)

	distribution, err := tc.service.CreateTaskDistribution(&req, userModel.ID)
	if err != nil {
		response.Fail(c, nil, "Failed to create task distribution: "+err.Error())
		return
	}

	response.Success(c, gin.H{"distribution": distribution}, "Task distribution created successfully")
}

// GetTaskDistributions handles retrieving all task distributions with pagination
func (tc *TaskDistributionController) GetTaskDistributions(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	distributions, total, err := tc.service.GetTaskDistributions(page, pageSize)
	if err != nil {
		response.Fail(c, nil, "Failed to retrieve task distributions: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"distributions": distributions,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	}, "Task distributions retrieved successfully")
}

// GetTaskDistribution handles retrieving a specific task distribution by ID
func (tc *TaskDistributionController) GetTaskDistribution(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid distribution ID")
		return
	}

	distribution, err := tc.service.GetTaskDistributionByID(uint(id))
	if err != nil {
		response.Fail(c, nil, "Failed to retrieve task distribution: "+err.Error())
		return
	}

	response.Success(c, gin.H{"distribution": distribution}, "Task distribution retrieved successfully")
}

// UpdateTaskDistribution handles updating an existing task distribution
func (tc *TaskDistributionController) UpdateTaskDistribution(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid distribution ID")
		return
	}

	var req model.TaskDistributionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, nil, "Invalid request data: "+err.Error())
		return
	}

	distribution, err := tc.service.UpdateTaskDistribution(uint(id), &req)
	if err != nil {
		response.Fail(c, nil, "Failed to update task distribution: "+err.Error())
		return
	}

	response.Success(c, gin.H{"distribution": distribution}, "Task distribution updated successfully")
}

// DeleteTaskDistribution handles deleting a task distribution
func (tc *TaskDistributionController) DeleteTaskDistribution(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid distribution ID")
		return
	}

	err = tc.service.DeleteTaskDistribution(uint(id))
	if err != nil {
		response.Fail(c, nil, "Failed to delete task distribution: "+err.Error())
		return
	}

	response.Success(c, nil, "Task distribution deleted successfully")
}

// ActivateTaskDistribution handles activating a task distribution
func (tc *TaskDistributionController) ActivateTaskDistribution(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid distribution ID")
		return
	}

	distribution, err := tc.service.ActivateTaskDistribution(uint(id))
	if err != nil {
		response.Fail(c, nil, "Failed to activate task distribution: "+err.Error())
		return
	}

	response.Success(c, gin.H{"distribution": distribution}, "Task distribution activated successfully")
}

// GetTaskDistributionsByTaskItem handles retrieving distributions for a specific task item
func (tc *TaskDistributionController) GetTaskDistributionsByTaskItem(c *gin.Context) {
	taskItemIdStr := c.Param("taskItemId")
	taskItemId, err := strconv.ParseUint(taskItemIdStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid task item ID")
		return
	}

	distributions, err := tc.service.GetTaskDistributionsByTaskItem(uint(taskItemId))
	if err != nil {
		response.Fail(c, nil, "Failed to retrieve task distributions: "+err.Error())
		return
	}

	response.Success(c, gin.H{"distributions": distributions}, "Task distributions retrieved successfully")
}

// GetDistributionSummary handles retrieving distribution summary statistics
func (tc *TaskDistributionController) GetDistributionSummary(c *gin.Context) {
	summary, err := tc.service.GetDistributionSummary()
	if err != nil {
		response.Fail(c, nil, "Failed to retrieve distribution summary: "+err.Error())
		return
	}

	response.Success(c, gin.H{"summary": summary}, "Distribution summary retrieved successfully")
}
