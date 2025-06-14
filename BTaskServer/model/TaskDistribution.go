package model

import (
	"fmt"
	"time"
)

// TaskDistribution represents a batch distribution of tasks to different user types
type TaskDistribution struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BatchName        string    `gorm:"type:varchar(255);not null" json:"batchName"`
	TaskItemID       uint      `gorm:"not null" json:"taskItemId"`
	TaskItem         TaskItem  `gorm:"foreignKey:TaskItemID" json:"taskItem"`
	TotalTasks       int       `gorm:"not null;default:0" json:"totalTasks"`
	RealMachineTasks int       `gorm:"default:0" json:"realMachineTasks"`
	ProtocolTasks    int       `gorm:"default:0" json:"protocolTasks"`
	ManualTasks      int       `gorm:"default:0" json:"manualTasks"`
	DistributedTasks int       `gorm:"default:0" json:"distributedTasks"`                // Sum of all assigned tasks
	RemainingTasks   int       `gorm:"default:0" json:"remainingTasks"`                  // TotalTasks - DistributedTasks
	Status           string    `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, active, completed, cancelled
	CreatedBy        uint      `gorm:"default:0" json:"createdBy"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName specifies the table name for TaskDistribution
func (TaskDistribution) TableName() string {
	return "task_distribution"
}

// TaskDistributionRequest represents the request payload for creating a task distribution
type TaskDistributionRequest struct {
	BatchName        string `json:"batchName" binding:"required"`
	TaskItemID       uint   `json:"taskItemId" binding:"required"`
	TotalTasks       int    `json:"totalTasks" binding:"required,min=1"`
	RealMachineTasks int    `json:"realMachineTasks" binding:"min=0"`
	ProtocolTasks    int    `json:"protocolTasks" binding:"min=0"`
	ManualTasks      int    `json:"manualTasks" binding:"min=0"`
}

// TaskDistributionUpdateRequest represents the request payload for updating a task distribution
type TaskDistributionUpdateRequest struct {
	BatchName        string `json:"batchName"`
	TotalTasks       int    `json:"totalTasks" binding:"min=1"`
	RealMachineTasks int    `json:"realMachineTasks" binding:"min=0"`
	ProtocolTasks    int    `json:"protocolTasks" binding:"min=0"`
	ManualTasks      int    `json:"manualTasks" binding:"min=0"`
	Status           string `json:"status"`
}

// Validate validates the task distribution request
func (req *TaskDistributionRequest) Validate() error {
	distributedSum := req.RealMachineTasks + req.ProtocolTasks + req.ManualTasks
	if distributedSum > req.TotalTasks {
		return fmt.Errorf("distributed tasks (%d) cannot exceed total tasks (%d)", distributedSum, req.TotalTasks)
	}
	if distributedSum == 0 {
		return fmt.Errorf("at least one user type must have tasks assigned")
	}
	return nil
}

// CalculateDistribution calculates the distributed and remaining tasks
func (td *TaskDistribution) CalculateDistribution() {
	td.DistributedTasks = td.RealMachineTasks + td.ProtocolTasks + td.ManualTasks
	td.RemainingTasks = td.TotalTasks - td.DistributedTasks
}
