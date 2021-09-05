package health

import (
	"github.com/jinzhu/gorm"
)

func ConfigureHealthChecks(db *gorm.DB) {
	NewCheckerService().
		Add(NewDbHealthChecker("defaultConnection", db))
}
