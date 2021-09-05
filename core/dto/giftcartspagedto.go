package dto

import "giftcard-engine/utils/indraframework"

type GiftCardsPageDTO struct {
	Size       int                            `json:"size"`
	Page       int                            `json:"page"`
	GiftCards  []GiftCardDTO                  `json:"gift_cards"`
	TotalItems int                            `json:"total_items"`
	Error      *indraframework.IndraException `json:"error"`
}

func NewGiftCardsPageDTO(giftCards []GiftCardDTO, size, page, total int) *GiftCardsPageDTO {
	return &GiftCardsPageDTO{
		Size:       size,
		Page:       page + 1,
		GiftCards:  giftCards,
		TotalItems: total,
	}
}

func (a *GiftCardsPageDTO) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
