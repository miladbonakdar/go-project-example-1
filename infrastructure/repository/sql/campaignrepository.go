package sql

import (
	"giftcard-engine/core"
	"giftcard-engine/core/common"
	"giftcard-engine/core/dbmodel"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

type campaignRepository struct {
	DB *gorm.DB
}

func (r *campaignRepository) FindByID(id uint) (dbmodel.Campaign, error) {
	var campaign dbmodel.Campaign

	if r.DB.Find(&campaign, id).RecordNotFound() {
		return dbmodel.EmptyCampaign(), common.CampaignNotFound
	}
	return campaign, nil
}

func (r *campaignRepository) Store(campaign *dbmodel.Campaign) error {
	err := r.titleGuard(campaign.Title)
	if err != nil {
		return err
	}
	return r.DB.Save(campaign).Error
}

func (r *campaignRepository) Delete(campaign dbmodel.Campaign) error {
	db := r.DB.Delete(&campaign)
	if db.RecordNotFound() {
		return common.CampaignNotFound
	}
	return db.Error
}

func (r *campaignRepository) titleGuard(title string) error {
	var total int
	r.DB.Model(&dbmodel.Campaign{}).Where("Title = ?", title).Count(&total)
	if total > 0 {
		return common.DuplicatedCampaignTitle
	}
	return nil
}

func (r *campaignRepository) FindPage(size, number uint, search string) ([]dbmodel.Campaign, int) {
	data := make(chan []dbmodel.Campaign)
	go func(channel chan<- []dbmodel.Campaign) {
		var campaigns []dbmodel.Campaign
		if search == "" {
			r.DB.Order("id desc").Limit(size).Offset(size * number).Find(&campaigns)
		} else {
			r.DB.Order("id desc").Where("Title like ?", "%"+search+"%").
				Limit(size).Offset(size * number).Find(&campaigns)
		}
		channel <- campaigns
	}(data)

	var total int
	if search == "" {
		r.DB.Model(&dbmodel.Campaign{}).Count(&total)
	} else {
		r.DB.Model(&dbmodel.Campaign{}).Where("Title like ?", "%"+search+"%").Count(&total)
	}
	return <-data, total
}

func NewCampaignRepository(DB *gorm.DB) core.CampaignRepository {
	return &campaignRepository{DB: DB}
}
