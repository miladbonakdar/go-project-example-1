package core

import (
	"giftcard-engine/core/dbmodel"
	"giftcard-engine/core/dto"
)

// Mapper is an interface to convert dto to sql model and vise versa
type Mapper interface {
	ToGiftCard(dto dto.CreateGiftCardDTO) *dbmodel.GiftCard
	ToGiftCardDTO(card *dbmodel.GiftCard) dto.GiftCardDTO
	ToListOfGiftCardDTO(cards []dbmodel.GiftCard) *dto.GiftCardsListDTO
	ToGiftCardStatusDTO(card dbmodel.GiftCard) dto.GiftCardStatusDTO
	ApprovedToGiftCardStatusDTO(card dbmodel.GiftCard) dto.GiftCardStatusDTO
	ToCampaign(dto dto.CreateCampaignDTO) dbmodel.Campaign
	ToCampaignDTO(campaign dbmodel.Campaign) dto.CampaignDTO
	ToListOfCampaigns(campaigns []dbmodel.Campaign) []dto.CampaignDTO
}
