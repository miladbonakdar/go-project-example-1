package dto

type BulkCreateGiftCardsDTO struct {
	GiftCards []CreateGiftCardDTO            `json:"gift_cards"`
}

func (a BulkCreateGiftCardsDTO) Validate() error {
	for _, card := range a.GiftCards {
		if err := card.Validate(); err != nil {
			return err
		}
	}
	return nil
}
