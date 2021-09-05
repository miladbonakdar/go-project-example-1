package api

import (
	"giftcard-engine/application/api/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"net/http"
)

func CreateRoute(cardHandler handlers.GiftCardHandler, campaignHandler handlers.CampaignHandler) *gin.Engine {
	route := gin.Default()
	giftCardV1 := route.Group("v1/gift-card")
	{
		giftCardV1.GET("/find/:id", cardHandler.FindByID)
		giftCardV1.POST("/", cardHandler.Store)
		giftCardV1.DELETE("/:id", cardHandler.Delete)
		giftCardV1.PUT("/", cardHandler.Update)
		giftCardV1.GET("/page/:size/:number", cardHandler.FindPage)
		giftCardV1.POST("/create-same-many", cardHandler.CreateSameMany)
		giftCardV1.POST("/create-many", cardHandler.CreateMany)
		giftCardV1.GET("/find-by-public-key/:key", cardHandler.FindByPublicKey)

		giftCardV1.POST("/approve-gift-cards", cardHandler.ApproveGiftCards)
		giftCardV1.POST("/validate-gift-cards", cardHandler.ValidateGiftCards)
		giftCardV1.PUT("/approve-gift-card/:uun/:secret", cardHandler.ApproveGiftCard)
		giftCardV1.GET("/validate-gift-card/:secret", cardHandler.ValidateGiftCard)
		giftCardV1.GET("/user-gift-cards/:uun", cardHandler.FindByUUN)

		giftCardV1.GET("/health", cardHandler.HealthCheck)
		giftCardV1.GET("/info", cardHandler.Info)
	}

	campaignV1 := route.Group("v1/campaign")
	{
		campaignV1.POST("/", campaignHandler.Create)
		campaignV1.PUT("/", campaignHandler.Update)
		campaignV1.DELETE("/:id", campaignHandler.Delete)
		campaignV1.GET("/page/:size/:number", campaignHandler.FindPage)
	}

	swaggerRedirectHandler := func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	}
	route.GET("/swagger", swaggerRedirectHandler)
	route.GET("/swagg", swaggerRedirectHandler)

	return route
}
