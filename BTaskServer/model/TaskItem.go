package model

import (
	"time"
)

type TaskItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TaskID        string    `gorm:"unique;not null" json:"taskId"`
	TaskName      string    `gorm:"not null" json:"taskName"`
	Supplier      string    `json:"supplier"`
	UserTypes     string    `json:"userTypes"` // Consider using a more structured type if complex logic is needed
	PriceReal     float64   `json:"priceReal"`
	PriceProtocol float64   `json:"priceProtocol"`
	PriceManual   float64   `json:"priceManual"`
	SelectionMode string    `json:"selectionMode"`
	Status        string    `json:"status" gorm:"default:'active'"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
