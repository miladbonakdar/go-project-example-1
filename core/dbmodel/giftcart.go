package dbmodel

import (
	"giftcard-engine/core/common"
	"giftcard-engine/utils/random"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mssql"
)

// GiftCard is a sql model for saving and modifying gift cards
type GiftCard struct {
	AbstractModel
	Amount     int32     `gorm:"column:Amount;not null"`
	PublicCode string    `gorm:"column:PublicCode;unique_index;not null"`
	SecretCode string    `gorm:"column:SecretCode;unique_index;not null"`
	UUN        string    `gorm:"column:UUN"`
	ExpireDate time.Time `gorm:"column:ExpireDate;not null"`
	Status     int       `gorm:"column:Status;not null;default:1"`
	CampaignId uint      `gorm:"column:CampaignId;not null;"`
	Campaign   *Campaign  `gorm:"foreignkey:ID;references:CampaignId"`
}

//TableName returns the sql table name for changing the default naming system
func (*GiftCard) TableName() string {
	return "GiftCard"
}

func NewGiftCard(amount int32, expireDate time.Time) *GiftCard {
	return &GiftCard{
		Amount:     amount,
		PublicCode: random.GiftCardPublicKey(),
		SecretCode: random.GiftCardSecretKey(),
		UUN:        "",
		ExpireDate: expireDate,
		Status:     Empty,
	}
}

func (g GiftCard) IsValid() bool {
	return g.UUN == "" && g.IsDateValid() && g.Status == Empty
}

func (g *GiftCard) SetCampaign(campaignId uint) {
	g.CampaignId = campaignId
}

func (g GiftCard) IsDateValid() bool {
	return time.Now().AddDate(0, 0, -1).UTC().Before(g.ExpireDate)
}

func (g *GiftCard) SetUUN(uun string) error {
	if g.UUN != "" {
		return common.GiftCardIsTaken
	}
	g.Status = Approved
	g.UUN = uun
	return nil
}

func (g *GiftCard) Update(amount int32, expireDate time.Time) error {
	if !g.IsValid() {
		return common.GiftCardIsNotValid
	}
	g.Amount = amount
	g.ExpireDate = expireDate
	return nil
}

func (g *GiftCard) RollBack() {
	g.UUN = ""
	g.Status = Empty
}

func (g *GiftCard) GenerateKey() {
	g.PublicCode = random.GiftCardPublicKey()
	g.SecretCode = random.GiftCardSecretKey()
}
