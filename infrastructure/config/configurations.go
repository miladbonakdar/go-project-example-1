package config

import (
	"giftcard-engine/infrastructure/config/configuration"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

var (
	once         sync.Once
	instance     configuration.Configurations
	instantiated bool
)

func readConfig() configuration.Configurations {
	once.Do(func() {
		err := godotenv.Load("dev.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		instance = configuration.CreateDefaultConfiguration()
		instantiated = true
	})

	return instance
}

func Get() configuration.Configurations {
	if !instantiated {
		return readConfig()
	}
	return instance
}
