package migrations

import (
	"github.com/spectrumwebco/agent_runtime/internal/database/models"
	"gorm.io/gorm"
)

// Migrate runs all migrations
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Agent{},
		&models.Tool{},
	)
}
