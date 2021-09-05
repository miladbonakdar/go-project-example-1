package dto

import "giftcard-engine/utils/indraframework"

// GiftCardDTO is an structure to get api input in gift card api.
type GiftCardDTO struct {
	ID            int                            `json:"id,string,omitempty"`
	PublicCode    string                         `json:"public_code"`
	SecretCode    string                         `json:"secret_code"`
	UUN           string                         `json:"uun"`
	ExpireDate    string                         `json:"expire_date"`
	Amount        int32                          `json:"amount"`
	IsValid       bool                           `json:"is_valid"`
	Error         *indraframework.IndraException `json:"error"`
	CampaignId    uint                           `json:"campaign_id"`
	CampaignTitle string                         `json:"campaign_title"`
}

func (a *GiftCardDTO) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
