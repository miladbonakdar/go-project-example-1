package core

import (
	"giftcard-engine/core/dto"
	"time"
)

// GiftCardService works with requests to api
type GiftCardService interface {
	FindPage(size, page uint, search string, campaignId *int,
		isValid *bool, expireDateFrom *time.Time, expireDateTo *time.Time) dto.GiftCardsPageDTO
	FindByID(id uint) (*dto.GiftCardDTO, error)
	Store(card *dto.CreateGiftCardDTO) (*dto.GiftCardDTO, error)
	Update(card *dto.UpdateGiftCardDto) (*dto.GiftCardDTO, error)
	Delete(id uint) error
	CreateMany(cards *dto.BulkCreateGiftCardsDTO) (*dto.GiftCardsListDTO, error)
	CreateSameMany(cards *dto.BulkCreateSameGiftCardsDTO) (*dto.GiftCardsListDTO, error)
	FindByPublicKey(key string) (*dto.GiftCardStatusDTO, error)

	FindByUUN(uun string) (*dto.GiftCardsListDTO, error)
	ValidateGiftCards(cards *dto.ValidateGiftCardsDto) *dto.GiftCardStatusListDTO
	ApproveGiftCards(cards *dto.ApproveGiftCardsDTO) (*dto.GiftCardStatusListDTO, error)
	ValidateGiftCard(giftCardSecret string) dto.GiftCardStatusDTO
	ApproveGiftCard(uun, giftCardSecret string) (dto.GiftCardStatusDTO, error)
}

type CampaignService interface {
	FindPage(size, page uint, search string) dto.CampaignPageDTO
	Create(campaign dto.CreateCampaignDTO) (dto.CampaignDTO, error)
	Update(campaign dto.UpdateCampaignDto) (dto.CampaignDTO, error)
	Delete(id uint) error
}
