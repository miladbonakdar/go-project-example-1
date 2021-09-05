package dbmodel

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mssql"
)

//Model Using int instead of uint for default gorm model
type AbstractModel struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
