package models

import (
	"time"
)

type Menu struct {
	ID         uint      `gorm:"primaryKey" json:"id" example:"1"`
	ParentID   *uint     `json:"parent_id"`
	Title      string    `gorm:"size:255;not null" json:"title" example:"Dashboard"`
	Path       *string   `gorm:"size:255" json:"path,omitempty" example:"/dashboard"`
	Icon       *string   `gorm:"size:100" json:"icon,omitempty" example:"icon-dashboard"`
	OrderIndex int       `gorm:"default:0" json:"order_index" example:"0"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Children   []Menu    `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}
