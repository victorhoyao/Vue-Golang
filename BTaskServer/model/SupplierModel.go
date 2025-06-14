package model

import (
	"time"
)

// Supplier represents the supplier model in the database
type Supplier struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SupplierID    string    `gorm:"type:varchar(255);not null;unique" json:"supplierId"`
	SupplierName  string    `gorm:"type:varchar(255);not null" json:"supplierName"`
	Contact       string    `gorm:"type:varchar(255)" json:"contact"`
	Status        string    `gorm:"type:varchar(50);not null;default:'Active'" json:"status"`
	TasksAssigned int       `gorm:"type:int;not null;default:0" json:"tasksAssigned"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// SupplierCreateForm represents the data for creating a new supplier
type SupplierCreateForm struct {
	SupplierID    string `json:"supplierId" binding:"required"`
	SupplierName  string `json:"supplierName" binding:"required"`
	Contact       string `json:"contact"`
	Status        string `json:"status"`
	TasksAssigned int    `json:"tasksAssigned"`
}

// SupplierUpdateForm represents the data for updating an existing supplier
type SupplierUpdateForm struct {
	SupplierID    string `json:"supplierId"`
	SupplierName  string `json:"supplierName"`
	Contact       string `json:"contact"`
	Status        string `json:"status"`
	TasksAssigned int    `json:"tasksAssigned"`
}
