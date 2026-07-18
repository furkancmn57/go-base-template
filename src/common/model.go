package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model is the base struct every GORM entity in the codebase must embed.
// It intentionally avoids relying on database-side UUID generation so the
// template works against any Postgres instance without extra extensions.
type Model struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate assigns a UUID before the row is inserted, unless one has
// already been set by the caller.
func (m *Model) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
