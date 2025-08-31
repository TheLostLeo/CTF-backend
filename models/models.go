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

// Challenge represents a CTF challenge
type Challenge struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Title       string         `json:"title" gorm:"not null" binding:"required"`
	Description string         `json:"description" gorm:"type:text"`
	Category    string         `json:"category" gorm:"not null" binding:"required"`
	Points      int            `json:"points" gorm:"not null" binding:"required"`
	Flag        string         `json:"-" gorm:"not null"` // Hidden from JSON
	Hint        string         `json:"hint"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete support

	// File attachments (optional)
	FileURL string `json:"file_url,omitempty"`

	// Relationships
	Submissions []Submission `json:"submissions,omitempty" gorm:"foreignKey:ChallengeID"`
}

// Submission represents a flag submission
type Submission struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	ChallengeID uint           `json:"challenge_id" gorm:"not null"`
	Flag        string         `json:"flag" gorm:"not null"`
	IsCorrect   bool           `json:"is_correct" gorm:"default:false"`
	IPAddress   string         `json:"ip_address,omitempty"`
	SubmittedAt time.Time      `json:"submitted_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete support

	// Relations
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Challenge Challenge `json:"challenge,omitempty" gorm:"foreignKey:ChallengeID"`
}
