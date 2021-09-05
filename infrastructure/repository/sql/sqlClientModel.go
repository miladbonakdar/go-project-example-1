package sql

import (
	"giftcard-engine/core/dbmodel"
	"giftcard-engine/infrastructure/logger"
	"github.com/jinzhu/gorm"
)

func newSqlClient(connectionString string) *gorm.DB {
	db, err := gorm.Open("mssql", connectionString)
	if err != nil {
		logger.FatalException(err, "Error creating connection pool")
	}

	db.SetLogger(&logger.GormLogger{})
	logger.Print("Connected!\n")
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(10)
	db.AutoMigrate(&dbmodel.GiftCard{}, &dbmodel.Campaign{})
	return db
}

// InitDatabase returns an implementation of the sql database orm with given connection string.
func InitDatabase(connectionString string) *gorm.DB {
	db := newSqlClient(connectionString)
	return db
}
