package handlers

import (
	"giftcard-engine/core"
	"giftcard-engine/core/common"
	"giftcard-engine/core/dto"
	"giftcard-engine/infrastructure/health"
	"giftcard-engine/utils"
	"giftcard-engine/utils/date"
	_ "giftcard-engine/utils/indraframework"
	"giftcard-engine/utils/parser"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type GiftCardHandler interface {
	FindByID(c *gin.Context)
	FindPage(c *gin.Context)
	Store(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	CreateMany(c *gin.Context)
	CreateSameMany(c *gin.Context)
	FindByPublicKey(c *gin.Context)

	ValidateGiftCards(c *gin.Context)
	ApproveGiftCards(c *gin.Context)
	ValidateGiftCard(c *gin.Context)
	ApproveGiftCard(c *gin.Context)
	FindByUUN(c *gin.Context)
	HealthCheck(c *gin.Context)
	Info(c *gin.Context)
}

type cardHandler struct {
	service core.GiftCardService
}

// FindByID godoc
// @Summary Get an gift card by id
// @Description find a gift card from the db
// @ID get-gift-card-by-id
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param id path int true "Gift Card ID"
// @Success 200 {object} dto.GiftCardDTO
// @Failure 400 {object} indraframework.IndraException
// @Failure 404 {object} indraframework.IndraException
// @Router /v1/gift-card/find/{id} [get]
func (h *cardHandler) FindByID(c *gin.Context) {
	id, err := parser.ParseNumber(c.Param("id"))
	if err != nil {
		jsonBadRequest(c, &dto.GiftCardDTO{}, err)
		return
	}
	giftCard, err := h.service.FindByID(id)

	if err == common.GiftCardNotFound {
		jsonNotFound(c, &dto.GiftCardDTO{}, err)
		return
	}

	jsonSuccess(c, giftCard)
}

// FindByPublicKey godoc
// @Summary Get gift card details by public key
// @Description find a gift card from the db
// @ID find-by-public-key
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param key path string true "Gift Card public key"
// @Success 200 {object} dto.GiftCardStatusDTO
// @Failure 404 {object} indraframework.IndraException
// @Router /v1/gift-card/find-by-public-key/{key} [get]
func (h *cardHandler) FindByPublicKey(c *gin.Context) {
	key := strings.ToUpper(c.Param("key"))
	giftCard, err := h.service.FindByPublicKey(key)

	if err == common.GiftCardNotFound {
		jsonNotFound(c, &dto.GiftCardStatusDTO{}, err)
		return
	}

	jsonSuccess(c, giftCard)
}

// FindByPublicKey godoc
// @Summary store a gift card
// @Description store a new gift card and generates the keys
// @ID store
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param giftCardDto body dto.CreateGiftCardDTO true "Create Gift Card dto"
// @Success 200 {object} dto.GiftCardDTO
// @Failure 400 {object} indraframework.IndraException
// @Failure 500 {object} indraframework.IndraException
// @Router /v1/gift-card [post]
func (h *cardHandler) Store(c *gin.Context) {
	var giftCardDTO dto.CreateGiftCardDTO
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&giftCardDTO), &dto.GiftCardDTO{} },
		func() (error error, data dto.Dto) { return giftCardDTO.Validate(), &dto.GiftCardDTO{} }); !success {
		return
	}

	giftCard, err := h.service.Store(&giftCardDTO)
	if err == common.InvalidCampaign {
		jsonBadRequest(c, &dto.GiftCardDTO{}, err)
	} else if err != nil {
		jsonInternalServerError(c, &dto.GiftCardDTO{}, err)
	} else {
		jsonSuccess(c, giftCard)
	}
}

// Update godoc
// @Summary updates a gift card
// @Description updates a gift card. just the expire date and the amount can be updated
// @ID update
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param giftCardDTO body dto.UpdateGiftCardDto true "Update Gift Card dto"
// @Success 200 {object} dto.GiftCardDTO
// @Failure 400 {object} indraframework.IndraException
// @Failure 404 {object} indraframework.IndraException
// @Failure 500 {object} indraframework.IndraException
// @Router /v1/gift-card [put]
func (h *cardHandler) Update(c *gin.Context) {
	var giftCardDTO dto.UpdateGiftCardDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&giftCardDTO), &dto.GiftCardDTO{} },
		func() (error error, data dto.Dto) { return giftCardDTO.Validate(), &dto.GiftCardDTO{} }); !success {
		return
	}

	giftCard, err := h.service.Update(&giftCardDTO)

	if err == common.GiftCardNotFound {
		jsonNotFound(c, &dto.GiftCardDTO{}, err)
	} else if err == common.GiftCardIsNotValid {
		jsonBadRequest(c, &dto.GiftCardDTO{}, err)
	} else if err != nil {
		jsonInternalServerError(c, &dto.GiftCardDTO{}, err)
	} else {
		jsonSuccess(c, giftCard)
	}

}

// Update godoc
// @Summary deletes a gift card
// @Description deletes a gift card by id
// @ID delete
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param id path int true "Gift Card's id"
// @Success 200 {object} dto.DeleteMessageDTO
// @Failure 400 {object} indraframework.IndraException
// @Failure 404 {object} indraframework.IndraException
// @Router /v1/gift-card/{id} [delete]
func (h *cardHandler) Delete(c *gin.Context) {
	id, err := parser.ParseNumber(c.Param("id"))
	if err != nil {
		jsonBadRequest(c, &dto.DeleteMessageDTO{}, err)
		return
	}

	err = h.service.Delete(id)

	if err == common.GiftCardNotFound {
		jsonNotFound(c, &dto.DeleteMessageDTO{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.DeleteMessageDTO{}, err)
		return
	}
	jsonSuccess(c, &dto.DeleteMessageDTO{
		ID:      int(id),
		Message: "Gift Card Deleted!",
		Error:   nil,
	})
}

// CreateMany godoc
// @Summary bulk insert gift cards
// @Description bulk insert for different gift cards
// @ID create-many
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param createGiftCards body dto.BulkCreateGiftCardsDTO true "bulk insert gift cards list"
// @Success 200 {object} dto.GiftCardsListDTO
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/gift-card/create-many [post]
func (h *cardHandler) CreateMany(c *gin.Context) {
	var createGiftCards dto.BulkCreateGiftCardsDTO
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&createGiftCards), &dto.GiftCardsListDTO{} },
		func() (error error, data dto.Dto) { return createGiftCards.Validate(), &dto.GiftCardsListDTO{} }); !success {
		return
	}
	cards, _ := h.service.CreateMany(&createGiftCards)
	jsonSuccess(c, cards)
}

// CreateSameMany godoc
// @Summary bulk insert gift cards
// @Description bulk insert for the same gift cards
// @ID create-same-many
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param createGiftCards body dto.BulkCreateSameGiftCardsDTO true "bulk insert for the same gift cards dto"
// @Success 200 {object} dto.GiftCardsListDTO
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/gift-card/create-same-many [post]
func (h *cardHandler) CreateSameMany(c *gin.Context) {
	var createGiftCards dto.BulkCreateSameGiftCardsDTO

	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&createGiftCards), &dto.GiftCardsListDTO{} },
		func() (error error, data dto.Dto) { return createGiftCards.Validate(), &dto.GiftCardsListDTO{} }); !success {
		return
	}
	cards, _ := h.service.CreateSameMany(&createGiftCards)
	jsonSuccess(c, cards)
}

// CreateSameMany godoc
// @Summary bulk validate gift cards
// @Description bulk validate for gift cards
// @ID validate-gift-cards
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param validateGiftCardsDto body dto.ValidateGiftCardsDto true "bulk validate dto"
// @Success 200 {object} dto.GiftCardStatusListDTO
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/gift-card/validate-gift-cards [post]
func (h *cardHandler) ValidateGiftCards(c *gin.Context) {
	var validateGiftCardsDto dto.ValidateGiftCardsDto

	if success := tryActions(c,
		func() (error error, data dto.Dto) {
			return c.BindJSON(&validateGiftCardsDto), &dto.GiftCardStatusListDTO{}
		},
		func() (error error, data dto.Dto) {
			return validateGiftCardsDto.Validate(), &dto.GiftCardStatusListDTO{}
		}); !success {
		return
	}
	cardsStatus := h.service.ValidateGiftCards(&validateGiftCardsDto)
	jsonSuccess(c, cardsStatus)
}

// ApproveGiftCards godoc
// @Summary bulk approve gift cards
// @Description bulk approve for gift cards
// @ID approve-gift-cards
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param approveGiftCardsDto body dto.ApproveGiftCardsDTO true "bulk approve dto"
// @Success 200 {object} dto.GiftCardStatusListDTO
// @Failure 400 {object} indraframework.IndraException
// @Failure 500 {object} indraframework.IndraException
// @Router /v1/gift-card/approve-gift-cards [post]
func (h *cardHandler) ApproveGiftCards(c *gin.Context) {
	var approveGiftCardsDto dto.ApproveGiftCardsDTO

	if success := tryActions(c,
		func() (error error, data dto.Dto) {
			return c.BindJSON(&approveGiftCardsDto), &dto.GiftCardStatusListDTO{}
		},
		func() (error error, data dto.Dto) {
			return approveGiftCardsDto.Validate(), &dto.GiftCardStatusListDTO{}
		}); !success {
		return
	}
	approveGiftCards, err := h.service.ApproveGiftCards(&approveGiftCardsDto)
	if err == common.GiftCardIsTaken {
		jsonBadRequest(c, &dto.GiftCardStatusListDTO{}, err)
		return
	}
	if err == common.GiftCardNotFound {
		jsonNotFound(c, &dto.GiftCardStatusListDTO{}, err)
		return
	}
	if err != nil {
		jsonInternalServerError(c, &dto.GiftCardStatusListDTO{}, err)
		return
	}
	jsonSuccess(c, approveGiftCards)
}

// FindPage godoc
// @Summary gift cards paging
// @Description get list of gift cards in paging object
// @ID find-page
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param size path integer true "page size"
// @Param number path integer true "page number"
// @Param campaignId query integer false "campaign id"
// @Param search query string false "search in public key"
// @Param isValid query boolean false "is valid gift card"
// @Param expireDateFrom query string false "expire date from"
// @Param expireDateTo query string false "expire date to"
// @Success 200 {object} dto.GiftCardsPageDTO
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/gift-card/page/{size}/{number} [get]
func (h *cardHandler) FindPage(c *gin.Context) {
	number, err := parser.ParseNumber(c.Param("number"))
	if err != nil {
		jsonBadRequest(c, &dto.GiftCardsPageDTO{}, err)
		return
	}
	var size uint
	size, err = parser.ParseNumber(c.Param("size"))
	if err != nil {
		jsonBadRequest(c, &dto.GiftCardsPageDTO{}, err)
		return
	}
	size = utils.MinUint(size, 50)
	if number == 0 {
		number += 1
	}
	number = number - 1
	var campaignId *int
	campaignIdQuery := c.Query("campaignId")
	if campaignIdQuery != "" {
		id, err := strconv.Atoi(campaignIdQuery)
		if err != nil {
			jsonBadRequest(c, &dto.GiftCardsPageDTO{}, common.InvalidCampaignQueryParam)
			return
		}
		campaignId = &id
	}
	search := c.Query("search")

	var isValid *bool
	isValidQuery := c.Query("isValid")
	if isValidQuery != "" {
		b, err := strconv.ParseBool(isValidQuery)
		if err != nil {
			jsonBadRequest(c, &dto.GiftCardsPageDTO{}, err)
			return
		}
		isValid = &b
	}

	var expireDateFrom *time.Time
	expireDateFromQuery := c.Query("expireDateFrom")
	if expireDateFromQuery != "" {
		d, err := date.DefaultToTime(expireDateFromQuery)
		if err != nil {
			jsonBadRequest(c, &dto.GiftCardsPageDTO{}, err)
			return
		}
		expireDateFrom = &d
	}

	var expireDateTo *time.Time
	expireDateToQuery := c.Query("expireDateTo")
	if expireDateToQuery != "" {
		d, err := date.DefaultToTime(expireDateToQuery)
		if err != nil {
			jsonBadRequest(c, &dto.GiftCardsPageDTO{}, err)
			return
		}
		expireDateTo = &d
	}

	cardsPage := h.service.FindPage(size, number, search, campaignId, isValid, expireDateFrom, expireDateTo)
	jsonSuccess(c, cardsPage)
}

// ValidateGiftCard godoc
// @Summary validate gift card
// @Description validate gift card
// @ID validate-gift-card
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param secret path string true "gift card secret"
// @Success 200 {object} dto.GiftCardStatusDTO
// @Router /v1/gift-card/validate-gift-card/{secret} [get]
func (h *cardHandler) ValidateGiftCard(c *gin.Context) {
	secret := c.Param("secret")
	status := h.service.ValidateGiftCard(secret)
	jsonSuccess(c, status)
}

// ApproveGiftCard godoc
// @Summary approve gift card
// @Description approve for single gift card
// @ID approve-gift-card
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param uun path string true "uun"
// @Param secret path string true "gift card secret"
// @Success 200 {object} dto.GiftCardStatusDTO
// @Failure 404 {object} indraframework.IndraException
// @Failure 400 {object} indraframework.IndraException
// @Failure 500 {object} indraframework.IndraException
// @Router /v1/gift-card/approve-gift-card/{uun}/{secret} [put]
func (h *cardHandler) ApproveGiftCard(c *gin.Context) {
	secret := c.Param("secret")
	uun := c.Param("uun")

	card, err := h.service.ApproveGiftCard(uun, secret)
	if err == common.GiftCardNotFound {
		jsonNotFound(c, &dto.GiftCardStatusDTO{}, err)
		return
	}

	if err == common.GiftCardIsTaken {
		jsonBadRequest(c, &dto.GiftCardStatusDTO{}, err)
		return
	}

	if err != nil {
		jsonInternalServerError(c, &dto.GiftCardStatusDTO{}, err)
		return
	}
	jsonSuccess(c, card)
}

// FindByUUN godoc
// @Summary user gift cards
// @Description get list of user's gift cards
// @ID find-by-uun
// @Accept  json
// @Produce  json
// @tags Gift Card
// @Param uun path string true "uun"
// @Success 200 {object} dto.GiftCardsListDTO
// @Failure 404 {object} indraframework.IndraException
// @Router /v1/gift-card/user-gift-cards/{uun} [get]
func (h *cardHandler) FindByUUN(c *gin.Context) {
	uun := c.Param("uun")

	giftCards, err := h.service.FindByUUN(uun)
	if err == common.NoGiftCardFoundForUser {
		jsonNotFound(c, &dto.GiftCardsListDTO{}, err)
		return
	}

	jsonSuccess(c, giftCards)
}

// HealthCheck godoc
// @Summary test Endpoint
// @Description get 200 response
// @Accept  json
// @Produce  json
// @tags Public
// @Success 200
// @Router /v1/gift-card/health [get]
func (h *cardHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, health.GetHealth())
}

// Info godoc
// @Summary Info
// @Description get 200 response
// @Accept  json
// @Produce  json
// @tags Public
// @Success 200
// @Router /v1/gift-card/info [get]
func (h *cardHandler) Info(c *gin.Context) {
	jsonSuccess(c, "Gift card api")
}

func NewGiftCardHandler(service core.GiftCardService) GiftCardHandler {
	return &cardHandler{service: service}
}
