package dto

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type CreateGiftCardDTO struct {
	ExpireDate string `json:"expire_date"`
	Amount     int32  `json:"amount"`
	CampaignId uint   `json:"campaign_id"`
}

func (a CreateGiftCardDTO) Validate() error {
	if err := CheckForDate(a.ExpireDate); err != nil {
		return err
	}
	return validation.ValidateStruct(&a,
		validation.Field(&a.ExpireDate, validation.Required),
		validation.Field(&a.CampaignId, validation.Required),
		validation.Field(&a.Amount, validation.Required, validation.Min(int32(1000))),
	)
}
