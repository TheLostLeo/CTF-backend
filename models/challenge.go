package models

import (
	"time"

	"gorm.io/gorm"
)

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

// TableName overrides the table name used by Challenge to `challenges`
func (Challenge) TableName() string {
	return "challenges"
}
