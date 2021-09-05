package configuration

import (
	"log"
	"os"
	"strconv"
)

type Configurations struct {
	Server            ServerConfiguration
	ConnectionStrings DatabaseConfiguration
	ElasticUrl        string
	ElasticHost       string
	ServiceName       string
	Environment       string
}

func (l Configurations) IsProduction() bool {
	return l.Environment == "Production" || l.Environment == "production"
}

func (l Configurations) IsStaging() bool {
	return l.Environment == "Staging" || l.Environment == "staging"
}

func (l Configurations) IsDevelopment() bool {
	return l.Environment == "Development" || l.Environment == "development"
}

func CreateDefaultConfiguration() Configurations {
	port, err := strconv.Atoi(os.Getenv("GIFT_CARD_SERVER_PORT"))
	if err != nil {
		log.Fatalln("The port number is not valid")
	}

	outSideOfContainerPort, err := strconv.Atoi(os.Getenv("GIFT_CARD_CONTAINER_PORT"))
	if err != nil {
		log.Fatalln("The container port number is not valid")
	}

	return Configurations{
		Server: ServerConfiguration{
			Port:                   port,
			OutSideOfContainerPort: outSideOfContainerPort,
			OutSideOfContainerHost: os.Getenv("GIFT_CARD_CONTAINER_NAME"),
		},
		ConnectionStrings: DatabaseConfiguration{
			DefaultConnection: os.Getenv("ConnectionStrings__DefaultConnection"),
		},
		Environment: os.Getenv("GIFT_CARD_ENVIRONMENT"),
		ElasticHost: os.Getenv("GIFT_CARD_ELASTIC_HOST"),
		ElasticUrl:  os.Getenv("GIFT_CARD_ELASTIC_URL"),
		ServiceName: os.Getenv("APP_NAME"),
	}
}
