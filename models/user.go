package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a CTF participant
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"unique;not null" binding:"required"`
	Email     string         `json:"email" gorm:"unique;not null" binding:"required,email"`
	Password  string         `json:"-" gorm:"not null"` // Hidden from JSON
	Score     int            `json:"score" gorm:"default:0"`
	IsAdmin   bool           `json:"is_admin" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete support

	// Relationships
	Submissions []Submission `json:"submissions,omitempty" gorm:"foreignKey:UserID"`
}

// TableName overrides the table name used by User to `users`
func (User) TableName() string {
	return "users"
}
