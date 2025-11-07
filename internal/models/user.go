package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"type:varchar(255);not null" json:"name"`
	Email     string          `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string          `gorm:"type:varchar(255);not null" json:"-"`
	Role      string          `gorm:"type:varchar(50);default:'user'" json:"role"`
	IsActive  bool            `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// GetPublicUser returns user data without sensitive fields
func (u *User) GetPublicUser() map[string]interface{} {
	return map[string]interface{}{
		"id":        u.ID,
		"name":      u.Name,
		"email":     u.Email,
		"role":      u.Role,
		"is_active": u.IsActive,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}
