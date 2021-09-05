package health

import (
	"giftcard-engine/infrastructure/logger"
	"github.com/jinzhu/gorm"
)

type dbHealthChecker struct {
	DB  *gorm.DB
	tag string
}

func (c *dbHealthChecker) Check() HealthResultDto {
	if err := c.DB.Exec("select 1;").Error; err != nil {
		logger.WithException(err).Error("error while trying to check if the sql connection is okay")
		return HealthResultDto{
			Status:      UnHealthy,
			Duration:    defaultTimeStampFormat,
			Exception:   err.Error(),
			Description: err.Error(),
			Data:        map[string]string{},
		}
	}
	return HealthResultDto{
		Status:   Healthy,
		Duration: defaultTimeStampFormat,
		Data:     map[string]string{},
	}
}

func (c *dbHealthChecker) Tag() string {
	return c.tag
}

func NewDbHealthChecker(tag string, db *gorm.DB) Checker {
	return &dbHealthChecker{
		DB:  db,
		tag: tag,
	}
}
