package models

import (
	"time"

	"gorm.io/gorm"
)

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

// TableName overrides the table name used by Submission to `submissions`
func (Submission) TableName() string {
	return "submissions"
}
