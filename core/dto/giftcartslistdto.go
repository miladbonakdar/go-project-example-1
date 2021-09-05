package dto

import "giftcard-engine/utils/indraframework"

// GiftCardDTO is an structure to get api input in gift card api.
type GiftCardsListDTO struct {
	Cards []GiftCardDTO                  `json:"gift_cards"`
	Error *indraframework.IndraException `json:"error"`
}

func (a *GiftCardsListDTO) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
