package sql

import (
	"giftcard-engine/core"
	"giftcard-engine/core/common"
	"giftcard-engine/core/dbmodel"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"time"
)

type gCardRepository struct {
	DB *gorm.DB
}

func (r *gCardRepository) FindByUUN(uun string) []dbmodel.GiftCard {
	var giftCards []dbmodel.GiftCard
	r.DB.Preload("Campaign").Find(&giftCards, "UUN=?", uun)
	return giftCards
}

func (r *gCardRepository) FindByID(id uint) (*dbmodel.GiftCard, error) {
	var giftCard dbmodel.GiftCard

	if r.DB.Preload("Campaign").Find(&giftCard, id).RecordNotFound() {
		return nil, common.GiftCardNotFound
	}
	return &giftCard, nil
}

func (r *gCardRepository) Store(card *dbmodel.GiftCard) error {
	err := r.campaignGuard(card.CampaignId)
	if err != nil {
		return err
	}
	return r.DB.Save(&card).Error
}

func (r *gCardRepository) campaignGuard(cid uint) error {
	if r.DB.Find(&dbmodel.Campaign{}, cid).RecordNotFound() {
		return common.InvalidCampaign
	}
	return nil
}

func (r *gCardRepository) Delete(card dbmodel.GiftCard) error {
	db := r.DB.Delete(&card)
	if db.RecordNotFound() {
		return common.GiftCardNotFound
	}
	return db.Error
}

func (r *gCardRepository) FindByPublicKey(key string) (*dbmodel.GiftCard, error) {
	var giftCard dbmodel.GiftCard

	if r.DB.Where("PublicCode = ?", key).First(&giftCard).RecordNotFound() {
		return nil, common.GiftCardNotFound
	}
	return &giftCard, nil
}

func (r *gCardRepository) FindPage(size, number uint, search string, campaignId *int,
	isValid *bool, expireDateFrom *time.Time, expireDateTo *time.Time) ([]dbmodel.GiftCard, int) {
	data := make(chan []dbmodel.GiftCard)

	query := r.DB.Model(&dbmodel.GiftCard{})
	if search != "" {
		query = query.Where("PublicCode like ?", "%"+search+"%")
	}
	if campaignId != nil {
		query = query.Where("CampaignId = ?", *campaignId)
	}
	if isValid != nil && *isValid == true {
		query = query.Where("(UUN is null or UUN = '') and ExpireDate > GETDATE()")
	}
	if isValid != nil && *isValid == false {
		query = query.Where("UUN is not null and ExpireDate < GETDATE()")
	}

	if expireDateFrom != nil {
		query = query.Where("ExpireDate > ?", *expireDateFrom)
	}

	if expireDateTo != nil {
		query = query.Where("ExpireDate < ?", *expireDateTo)
	}

	go func(channel chan<- []dbmodel.GiftCard) {
		var giftCards []dbmodel.GiftCard
		query.Order("id desc").Limit(size).Offset(size * number).Find(&giftCards)
		channel <- giftCards
	}(data)

	var total int
	query.Count(&total)
	return <-data, total
}

func (r *gCardRepository) FindBySecretKey(secret string) (*dbmodel.GiftCard, error) {
	var giftCard dbmodel.GiftCard

	if r.DB.Where("SecretCode = ?", secret).First(&giftCard).RecordNotFound() {
		return nil, common.GiftCardNotFound
	}
	return &giftCard, nil
}

func (r *gCardRepository) RollBackApprove(secret string) error {
	var giftCard dbmodel.GiftCard

	if r.DB.Where("SecretCode = ?", secret).First(&giftCard).RecordNotFound() {
		return common.GiftCardNotFound
	}
	giftCard.RollBack()
	return r.DB.Save(giftCard).Error
}

func NewGiftCardRepository(DB *gorm.DB) core.GiftCardRepository {
	return &gCardRepository{DB: DB}
}
