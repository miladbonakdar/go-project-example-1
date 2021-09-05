package main

import (
	"fmt"
	"giftcard-engine/application/api"
	"giftcard-engine/application/api/handlers"
	"giftcard-engine/cmd/docs"
	"giftcard-engine/core/logic"
	"giftcard-engine/infrastructure/config"
	"giftcard-engine/infrastructure/health"
	"giftcard-engine/infrastructure/logger"
	"giftcard-engine/infrastructure/repository/sql"
	"github.com/jinzhu/gorm"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

var db *gorm.DB

func main() {
	configurations := config.Get()

	db = sql.InitDatabase(configurations.ConnectionStrings.DefaultConnection)
	defer db.Close()
	logger.ConfigureLogger(
		logger.LoggerConfiguration{
			ServiceName: configurations.ServiceName,
			Environment: configurations.Environment,
			ElasticUrl:  configurations.ElasticUrl,
		})
	health.ConfigureHealthChecks(db)
	gRepository := sql.NewGiftCardRepository(db)
	campaignRepository := sql.NewCampaignRepository(db)
	gMapper := sql.NewMapper()
	gService := logic.NewGiftCardService(gRepository, gMapper)
	campaignService := logic.NewCampaignService(campaignRepository, gMapper)
	gHandler := handlers.NewGiftCardHandler(gService)
	cHandler := handlers.NewCampaignHandler(campaignService)
	//routes
	route := api.CreateRoute(gHandler, cHandler)
	//swagger
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%v", configurations.Server.OutSideOfContainerHost,
		configurations.Server.OutSideOfContainerPort)
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := route.Run(fmt.Sprintf(":%d", configurations.Server.Port))
	if err != nil {
		logger.Panic(err.Error())
	}
}
