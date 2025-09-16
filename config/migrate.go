package config

import (
	"github.com/linn221/RequesterBackend/models"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	// Auto-migrate all models in dependency order
	err := db.AutoMigrate(
		&models.Program{},    // No dependencies
		&models.ImportJob{},  // No dependencies
		&models.Endpoint{},   // Depends on Program
		&models.MyRequest{},  // Depends on Program, ImportJob, Endpoint
		&models.Attachment{}, // Polymorphic - depends on all above
		&models.Note{},       // Polymorphic - depends on all above
	)
	if err != nil {
		panic("Error migrating tables: " + err.Error())
	}
}
