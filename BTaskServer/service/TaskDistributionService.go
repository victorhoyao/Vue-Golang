package service

import (
	"BTaskServer/dao"
	"BTaskServer/model"
	"fmt"
)

// TaskDistributionService handles business logic for task distributions
type TaskDistributionService struct{}

// NewTaskDistributionService creates a new TaskDistributionService
func NewTaskDistributionService() *TaskDistributionService {
	return &TaskDistributionService{}
}

// CreateTaskDistribution creates a new task distribution
func (s *TaskDistributionService) CreateTaskDistribution(req *model.TaskDistributionRequest, createdBy uint) (*model.TaskDistribution, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if task item exists
	_, err := dao.GetTaskItemByID(req.TaskItemID)
	if err != nil {
		return nil, fmt.Errorf("task item not found: %v", err)
	}

	// Create task distribution
	distribution := &model.TaskDistribution{
		BatchName:        req.BatchName,
		TaskItemID:       req.TaskItemID,
		TotalTasks:       req.TotalTasks,
		RealMachineTasks: req.RealMachineTasks,
		ProtocolTasks:    req.ProtocolTasks,
		ManualTasks:      req.ManualTasks,
		CreatedBy:        createdBy,
		Status:           "pending",
	}

	// Calculate distribution
	distribution.CalculateDistribution()

	// Save to database
	if err := dao.CreateTaskDistribution(distribution); err != nil {
		return nil, fmt.Errorf("failed to create task distribution: %v", err)
	}

	// Load with associations
	return dao.GetTaskDistributionByID(distribution.ID)
}

// GetTaskDistributions retrieves all task distributions with pagination
func (s *TaskDistributionService) GetTaskDistributions(page, pageSize int) ([]model.TaskDistribution, int64, error) {
	return dao.GetTaskDistributions(page, pageSize)
}

// GetTaskDistributionByID retrieves a task distribution by ID
func (s *TaskDistributionService) GetTaskDistributionByID(id uint) (*model.TaskDistribution, error) {
	return dao.GetTaskDistributionByID(id)
}

// UpdateTaskDistribution updates an existing task distribution
func (s *TaskDistributionService) UpdateTaskDistribution(id uint, req *model.TaskDistributionUpdateRequest) (*model.TaskDistribution, error) {
	// Get existing distribution
	distribution, err := dao.GetTaskDistributionByID(id)
	if err != nil {
		return nil, fmt.Errorf("task distribution not found: %v", err)
	}

	// Check if distribution can be updated (completed distributions cannot be updated)
	if distribution.Status == "completed" {
		return nil, fmt.Errorf("cannot update completed distribution")
	}

	// Validate the update request
	distributedSum := req.RealMachineTasks + req.ProtocolTasks + req.ManualTasks
	if distributedSum > req.TotalTasks {
		return nil, fmt.Errorf("distributed tasks (%d) cannot exceed total tasks (%d)", distributedSum, req.TotalTasks)
	}

	// Update fields
	if req.BatchName != "" {
		distribution.BatchName = req.BatchName
	}
	if req.TotalTasks > 0 {
		distribution.TotalTasks = req.TotalTasks
	}
	distribution.RealMachineTasks = req.RealMachineTasks
	distribution.ProtocolTasks = req.ProtocolTasks
	distribution.ManualTasks = req.ManualTasks
	if req.Status != "" {
		distribution.Status = req.Status
	}

	// Recalculate distribution
	distribution.CalculateDistribution()

	// Save to database
	if err := dao.UpdateTaskDistribution(distribution); err != nil {
		return nil, fmt.Errorf("failed to update task distribution: %v", err)
	}

	return dao.GetTaskDistributionByID(id)
}

// DeleteTaskDistribution deletes a task distribution
func (s *TaskDistributionService) DeleteTaskDistribution(id uint) error {
	// Get existing distribution
	distribution, err := dao.GetTaskDistributionByID(id)
	if err != nil {
		return fmt.Errorf("task distribution not found: %v", err)
	}

	// Check if distribution can be deleted (completed distributions cannot be deleted)
	if distribution.Status == "completed" {
		return fmt.Errorf("cannot delete completed distribution")
	}

	return dao.DeleteTaskDistribution(id)
}

// ActivateTaskDistribution activates a task distribution (makes it available for task assignment)
func (s *TaskDistributionService) ActivateTaskDistribution(id uint) (*model.TaskDistribution, error) {
	distribution, err := dao.GetTaskDistributionByID(id)
	if err != nil {
		return nil, fmt.Errorf("task distribution not found: %v", err)
	}

	if distribution.Status != "pending" {
		return nil, fmt.Errorf("can only activate pending distributions, current status: %s", distribution.Status)
	}

	distribution.Status = "active"
	if err := dao.UpdateTaskDistribution(distribution); err != nil {
		return nil, fmt.Errorf("failed to activate task distribution: %v", err)
	}

	return dao.GetTaskDistributionByID(id)
}

// GetTaskDistributionsByTaskItem retrieves all distributions for a specific task item
func (s *TaskDistributionService) GetTaskDistributionsByTaskItem(taskItemID uint) ([]model.TaskDistribution, error) {
	return dao.GetTaskDistributionsByTaskItem(taskItemID)
}

// GetDistributionSummary provides a summary of task distributions
func (s *TaskDistributionService) GetDistributionSummary() (map[string]interface{}, error) {
	summary, err := dao.GetTaskDistributionSummary()
	if err != nil {
		return nil, fmt.Errorf("failed to get distribution summary: %v", err)
	}
	return summary, nil
}
