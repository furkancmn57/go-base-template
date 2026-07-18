// Package extensions wires infrastructure at the composition root.
package extensions

import (
	"gorm.io/gorm"

	"github.com/furkancmn57/go-base-template/src/config"
	"github.com/furkancmn57/go-base-template/src/data"
)

// AddDatabase opens Postgres and runs AutoMigrate.
func AddDatabase(cfg config.Postgres) (*gorm.DB, error) {
	db, err := data.New(cfg)
	if err != nil {
		return nil, err
	}
	if err := data.Migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}
