package entities

import (
	"github.com/furkancmn57/go-base-template/src/common"
	"github.com/furkancmn57/go-base-template/src/data/mappings"
)

// Todo is the persisted todo item.
type Todo struct {
	common.Model
	Title       string `gorm:"column:title;type:varchar(255);not null"`
	Description string `gorm:"column:description;type:text"`
	Completed   bool   `gorm:"column:completed;not null;default:false"`
}

// TableName returns the physical table name from mappings.
func (Todo) TableName() string {
	return mappings.TodoTable
}
