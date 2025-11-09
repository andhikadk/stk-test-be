package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Menu struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ParentID   *uuid.UUID `gorm:"type:uuid" json:"parent_id,omitempty"`
	Title      string     `gorm:"size:255;not null" json:"title" example:"Dashboard"`
	Path       *string    `gorm:"size:255" json:"path,omitempty" example:"/dashboard"`
	Icon       *string    `gorm:"size:100" json:"icon,omitempty" example:"icon-dashboard"`
	OrderIndex int        `gorm:"default:0" json:"order_index" example:"0"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Children   []Menu     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
