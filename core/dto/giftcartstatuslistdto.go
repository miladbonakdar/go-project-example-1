package dto

import "giftcard-engine/utils/indraframework"

type GiftCardStatusListDTO struct {
	Cards []GiftCardStatusDTO            `json:"gift_cards_statuses"`
	Error *indraframework.IndraException `json:"error"`
}

func (a *GiftCardStatusListDTO) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
