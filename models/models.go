package models

import (
	"gorm.io/gorm"
)

// GetAllModels returns a slice of all model types for migration
func GetAllModels() []interface{} {
	return []interface{}{
		&User{},
		&Challenge{},
		&Submission{},
	}
}

// MigrateAll runs auto-migration for all models
func MigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(GetAllModels()...)
}
