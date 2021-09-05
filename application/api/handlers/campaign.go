package handlers

import (
	"giftcard-engine/core"
	"giftcard-engine/core/common"
	"giftcard-engine/core/dto"
	"giftcard-engine/utils"
	_ "giftcard-engine/utils/indraframework"
	"giftcard-engine/utils/parser"
	"github.com/gin-gonic/gin"
)

type CampaignHandler interface {
	FindPage(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Create(c *gin.Context)
}

type campaignHandler struct {
	service core.CampaignService
}

// Campaign FindPage godoc
// @Summary Campaign paging
// @Description get list of campaigns in paging object
// @ID find-page
// @Accept  json
// @tags Campaign
// @Produce  json
// @Param size path number true "page size"
// @Param number path number true "page number"
// @Param search query string false "search by title"
// @Success 200 {object} dto.CampaignPageDTO
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/campaign/page/{size}/{number} [get]
func (h *campaignHandler) FindPage(c *gin.Context) {
	number, err := parser.ParseNumber(c.Param("number"))
	if err != nil {
		jsonBadRequest(c, &dto.CampaignPageDTO{}, err)
		return
	}
	var size uint
	size, err = parser.ParseNumber(c.Param("size"))
	if err != nil {
		jsonBadRequest(c, &dto.CampaignPageDTO{}, err)
		return
	}
	size = utils.MinUint(size, 50)
	if number == 0 {
		number += 1
	}
	number = number - 1
	campaignsPage := h.service.FindPage(size, number, c.Query("search"))
	jsonSuccess(c, campaignsPage)
}

// create a new campaign godoc
// @Summary store a campaign
// @Description store a new campaign and generates the keys
// @ID store
// @Accept  json
// @tags Campaign
// @Produce  json
// @Param campaignDTO body dto.CreateCampaignDTO true "Create campaign dto"
// @Success 200 {object} dto.CampaignDTO
// @Failure 400 {object} indraframework.IndraException
// @Failure 500 {object} indraframework.IndraException
// @Router /v1/campaign [post]
func (h *campaignHandler) Create(c *gin.Context) {
	var campaignDTO dto.CreateCampaignDTO
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&campaignDTO), &dto.CampaignDTO{} },
		func() (error error, data dto.Dto) { return campaignDTO.Validate(), &dto.CampaignDTO{} }); !success {
		return
	}

	campaign, err := h.service.Create(campaignDTO)
	if err == common.DuplicatedCampaignTitle {
		jsonBadRequest(c, &dto.CampaignDTO{}, err)
	} else if err != nil {
		jsonInternalServerError(c, &dto.CampaignDTO{}, err)
	} else {
		jsonSuccess(c, campaign)
	}
}

// Update godoc
// @Summary updates a campaign
// @Description updates a campaign
// @ID update
// @Accept  json
// @tags Campaign
// @Produce  json
// @Param campaignDTO body dto.UpdateCampaignDto true "Update Campaign dto"
// @Success 200 {object} dto.CampaignDTO
// @Failure 400 {object} indraframework.IndraException
// @Failure 404 {object} indraframework.IndraException
// @Failure 500 {object} indraframework.IndraException
// @Router /v1/campaign [put]
func (h *campaignHandler) Update(c *gin.Context) {
	var campaignDTO dto.UpdateCampaignDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&campaignDTO), &dto.CampaignDTO{} },
		func() (error error, data dto.Dto) { return campaignDTO.Validate(), &dto.CampaignDTO{} }); !success {
		return
	}

	campaign, err := h.service.Update(campaignDTO)

	if err == common.CampaignNotFound {
		jsonNotFound(c, &dto.CampaignDTO{}, err)
	} else if err == common.DuplicatedCampaignTitle {
		jsonBadRequest(c, &dto.CampaignDTO{}, err)
	} else if err != nil {
		jsonInternalServerError(c, &dto.CampaignDTO{}, err)
	} else {
		jsonSuccess(c, campaign)
	}

}

// Update godoc
// @Summary deletes a campaign
// @Description deletes a campaign by id
// @ID delete
// @Accept  json
// @tags Campaign
// @Produce  json
// @Param id path int true "campaign's id"
// @Success 200 {object} dto.DeleteMessageDTO
// @Failure 400 {object} indraframework.IndraException
// @Failure 404 {object} indraframework.IndraException
// @Router /v1/campaign/{id} [delete]
func (h *campaignHandler) Delete(c *gin.Context) {
	id, err := parser.ParseNumber(c.Param("id"))
	if err != nil {
		jsonBadRequest(c, &dto.DeleteMessageDTO{}, err)
		return
	}

	err = h.service.Delete(id)

	if err == common.CampaignNotFound {
		jsonNotFound(c, &dto.DeleteMessageDTO{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.DeleteMessageDTO{}, err)
		return
	}
	jsonSuccess(c, &dto.DeleteMessageDTO{
		ID:      int(id),
		Message: "Campaign Deleted!",
		Error:   nil,
	})
}

func NewCampaignHandler(service core.CampaignService) CampaignHandler {
	return &campaignHandler{service: service}
}
