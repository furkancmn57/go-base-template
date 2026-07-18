package data

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/furkancmn57/go-base-template/src/data/entities"
)

// Migrate auto-migrates registered entities. Append new entities here.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entities.Todo{},
	); err != nil {
		return fmt.Errorf("data: failed to auto-migrate: %w", err)
	}
	return nil
}
