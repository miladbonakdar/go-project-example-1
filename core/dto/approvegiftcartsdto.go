package dto

import (
	"giftcard-engine/utils"
	"github.com/go-ozzo/ozzo-validation/v4"
)

type ApproveGiftCardsDTO struct {
	UUN             string                         `json:"uun"`
	GiftCardsSecret []string                       `json:"gift_cards_secret"`
}

func (a ApproveGiftCardsDTO) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.UUN, validation.Required),
		validation.Field(&a.GiftCardsSecret,
			validation.Each(validation.Length(utils.GiftCardSecretKeyLength, utils.GiftCardSecretKeyLength))),
	)
}
