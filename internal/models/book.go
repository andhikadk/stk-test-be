package models

import (
	"time"

	"gorm.io/gorm"
)

// Book represents a book in the library
type Book struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Title     string          `gorm:"type:varchar(255);not null" json:"title"`
	Author    string          `gorm:"type:varchar(255);not null" json:"author"`
	ISBN      string          `gorm:"type:varchar(20);uniqueIndex" json:"isbn"`
	Year      int             `gorm:"type:int" json:"year"`
	Pages     int             `gorm:"type:int" json:"pages"`
	Publisher string          `gorm:"type:varchar(255)" json:"publisher"`
	Description string        `gorm:"type:text" json:"description"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName specifies the table name for Book model
func (Book) TableName() string {
	return "books"
}

// IsValid validates the book data
func (b *Book) IsValid() bool {
	if b.Title == "" || b.Author == "" || b.ISBN == "" {
		return false
	}
	if b.Year < 1000 || b.Year > 9999 {
		return false
	}
	return true
}
