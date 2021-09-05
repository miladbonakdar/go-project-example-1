package core

import (
	"giftcard-engine/core/dbmodel"
	"time"
)

type GiftCardRepository interface {
	FindByUUN(uun string) []dbmodel.GiftCard
	FindByID(id uint) (*dbmodel.GiftCard, error)
	Store(card *dbmodel.GiftCard) error
	Delete(card dbmodel.GiftCard) error
	FindByPublicKey(key string) (*dbmodel.GiftCard, error)
	FindPage(size, number uint, search string, campaignId *int, isValid *bool,
		expireDateFrom *time.Time, expireDateTo *time.Time) ([]dbmodel.GiftCard, int)
	FindBySecretKey(secret string) (*dbmodel.GiftCard, error)
	RollBackApprove(secret string) error
}

type CampaignRepository interface {
	FindByID(id uint) (dbmodel.Campaign, error)
	Store(card *dbmodel.Campaign) error
	Delete(card dbmodel.Campaign) error
	FindPage(size, number uint, search string) ([]dbmodel.Campaign, int)
}
