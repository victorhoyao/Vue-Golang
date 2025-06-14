package service

import (
	"BTaskServer/global"
	"BTaskServer/model"
)

// CreateTaskItem creates a new task item in the database.
func CreateTaskItem(taskItem *model.TaskItem) (*model.TaskItem, error) {
	if err := global.GVA_DB.Create(taskItem).Error; err != nil {
		return nil, err
	}
	return taskItem, nil
}

// GetTaskItems retrieves all task items from the database.
func GetTaskItems() ([]model.TaskItem, error) {
	var taskItems []model.TaskItem
	if err := global.GVA_DB.Find(&taskItems).Error; err != nil {
		return nil, err
	}
	return taskItems, nil
}

// UpdateTaskItem updates an existing task item in the database.
func UpdateTaskItem(id uint, taskItem *model.TaskItem) (*model.TaskItem, error) {
	var existingTaskItem model.TaskItem
	if err := global.GVA_DB.Where("id = ?", id).First(&existingTaskItem).Error; err != nil {
		return nil, err
	}

	if err := global.GVA_DB.Model(&existingTaskItem).Updates(taskItem).Error; err != nil {
		return nil, err
	}

	// Fetch the updated task item to return it with all fields
	var updatedTaskItem model.TaskItem
	if err := global.GVA_DB.Where("id = ?", id).First(&updatedTaskItem).Error; err != nil {
		return nil, err
	}

	return &updatedTaskItem, nil
}

// DeleteTaskItem deletes a task item from the database.
func DeleteTaskItem(id uint) error {
	if err := global.GVA_DB.Where("id = ?", id).Delete(&model.TaskItem{}).Error; err != nil {
		return err
	}
	return nil
}
