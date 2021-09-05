package dto

import (
	"giftcard-engine/utils"
	"github.com/go-ozzo/ozzo-validation/v4"
)

type ValidateGiftCardsDto struct {
	GiftCardsSecret []string                       `json:"gift_cards_secret"`
}

func (a ValidateGiftCardsDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.GiftCardsSecret,
			validation.Each(validation.Length(utils.GiftCardSecretKeyLength, utils.GiftCardSecretKeyLength))),
	)
}
