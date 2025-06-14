package dao

import (
	"BTaskServer/global"
	"BTaskServer/model"
	"fmt"
)

// GetTaskItemByID retrieves a task item by ID (needed for task distribution validation)
func GetTaskItemByID(id uint) (*model.TaskItem, error) {
	var taskItem model.TaskItem
	err := global.GVA_DB.First(&taskItem, id).Error
	if err != nil {
		return nil, err
	}
	return &taskItem, nil
}

// CreateTaskDistribution creates a new task distribution in the database
func CreateTaskDistribution(distribution *model.TaskDistribution) error {
	return global.GVA_DB.Create(distribution).Error
}

// GetTaskDistributions retrieves all task distributions with pagination
func GetTaskDistributions(page, pageSize int) ([]model.TaskDistribution, int64, error) {
	var distributions []model.TaskDistribution
	var total int64

	// Count total records
	if err := global.GVA_DB.Model(&model.TaskDistribution{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated results with preloaded TaskItem
	err := global.GVA_DB.
		Preload("TaskItem").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&distributions).Error

	return distributions, total, err
}

// GetTaskDistributionByID retrieves a task distribution by ID
func GetTaskDistributionByID(id uint) (*model.TaskDistribution, error) {
	var distribution model.TaskDistribution
	err := global.GVA_DB.Preload("TaskItem").First(&distribution, id).Error

	if err != nil {
		return nil, err
	}
	return &distribution, nil
}

// UpdateTaskDistribution updates an existing task distribution
func UpdateTaskDistribution(distribution *model.TaskDistribution) error {
	return global.GVA_DB.Save(distribution).Error
}

// DeleteTaskDistribution deletes a task distribution by ID
func DeleteTaskDistribution(id uint) error {
	return global.GVA_DB.Delete(&model.TaskDistribution{}, id).Error
}

// GetTaskDistributionsByTaskItem retrieves all distributions for a specific task item
func GetTaskDistributionsByTaskItem(taskItemID uint) ([]model.TaskDistribution, error) {
	var distributions []model.TaskDistribution
	err := global.GVA_DB.
		Preload("TaskItem").
		Where("task_item_id = ?", taskItemID).
		Order("created_at DESC").
		Find(&distributions).Error

	return distributions, err
}

// GetTaskDistributionsByStatus retrieves distributions by status
func GetTaskDistributionsByStatus(status string) ([]model.TaskDistribution, error) {
	var distributions []model.TaskDistribution
	err := global.GVA_DB.
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&distributions).Error

	return distributions, err
}

// GetTaskDistributionSummary provides summary statistics for task distributions
func GetTaskDistributionSummary() (map[string]interface{}, error) {
	var summary map[string]interface{} = make(map[string]interface{})

	// Total distributions
	var totalDistributions int64
	if err := global.GVA_DB.Model(&model.TaskDistribution{}).Count(&totalDistributions).Error; err != nil {
		return nil, err
	}
	summary["totalDistributions"] = totalDistributions

	// Distributions by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	if err := global.GVA_DB.Model(&model.TaskDistribution{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, err
	}
	summary["statusCounts"] = statusCounts

	// Total tasks summary
	var taskSummary struct {
		TotalTasks       int64
		DistributedTasks int64
		RemainingTasks   int64
		RealMachineTasks int64
		ProtocolTasks    int64
		ManualTasks      int64
	}
	if err := global.GVA_DB.Model(&model.TaskDistribution{}).
		Select(`
			SUM(total_tasks) as total_tasks,
			SUM(distributed_tasks) as distributed_tasks,
			SUM(remaining_tasks) as remaining_tasks,
			SUM(real_machine_tasks) as real_machine_tasks,
			SUM(protocol_tasks) as protocol_tasks,
			SUM(manual_tasks) as manual_tasks
		`).
		Where("status IN ?", []string{"active", "completed"}).
		Scan(&taskSummary).Error; err != nil {
		return nil, err
	}
	summary["taskSummary"] = taskSummary

	// Recent distributions (last 10)
	var recentDistributions []model.TaskDistribution
	if err := global.GVA_DB.
		Order("created_at DESC").
		Limit(10).
		Find(&recentDistributions).Error; err != nil {
		return nil, err
	}
	summary["recentDistributions"] = recentDistributions

	return summary, nil
}

// GetActiveDistributionsForUserType retrieves active distributions that have available tasks for a specific user type
func GetActiveDistributionsForUserType(userType string) ([]model.TaskDistribution, error) {
	var distributions []model.TaskDistribution
	var whereClause string

	switch userType {
	case "Real Machine":
		whereClause = "status = 'active' AND real_machine_tasks > 0"
	case "Protocol":
		whereClause = "status = 'active' AND protocol_tasks > 0"
	case "Manual":
		whereClause = "status = 'active' AND manual_tasks > 0"
	default:
		return nil, fmt.Errorf("invalid user type: %s", userType)
	}

	err := global.GVA_DB.
		Preload("TaskItem").
		Where(whereClause).
		Order("created_at ASC").
		Find(&distributions).Error

	return distributions, err
}

// UpdateTaskDistributionTaskCounts updates the task counts for a distribution (used when tasks are assigned)
func UpdateTaskDistributionTaskCounts(id uint, userType string, assignedCount int) error {
	var updateField string

	switch userType {
	case "Real Machine":
		updateField = "real_machine_tasks"
	case "Protocol":
		updateField = "protocol_tasks"
	case "Manual":
		updateField = "manual_tasks"
	default:
		return fmt.Errorf("invalid user type: %s", userType)
	}

	// Update the specific user type task count and recalculate distributed/remaining tasks
	return global.GVA_DB.Exec(fmt.Sprintf(`
		UPDATE task_distributions 
		SET %s = %s - ?,
			distributed_tasks = real_machine_tasks + protocol_tasks + manual_tasks,
			remaining_tasks = total_tasks - (real_machine_tasks + protocol_tasks + manual_tasks),
			updated_at = NOW()
		WHERE id = ? AND %s >= ?
	`, updateField, updateField, updateField), assignedCount, id, assignedCount).Error
}
