package dto

import (
	"giftcard-engine/core/common"
	"giftcard-engine/utils/date"
	"github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type UpdateGiftCardDto struct {
	ID         int    `json:"id"`
	ExpireDate string `json:"expire_date"`
	Amount     int32  `json:"amount"`
}

func (a UpdateGiftCardDto) Validate() error {
	if err := CheckForDate(a.ExpireDate); err != nil {
		return err
	}

	return validation.ValidateStruct(&a,
		validation.Field(&a.ExpireDate, validation.Required),
		validation.Field(&a.ID, validation.Required),
		validation.Field(&a.Amount, validation.Required, validation.Min(int32(1000))),
	)
}

func CheckForDate(expireDate string) error {
	expDate, err := date.DefaultToTime(expireDate)
	if err != nil {
		return err
	}
	if expDate.Before(time.Now().AddDate(0, 0, -1).UTC()) {
		return common.ExpireDateIsNotInValidRange
	}
	return nil
}
