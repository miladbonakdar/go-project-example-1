package dto

import "giftcard-engine/utils/indraframework"

// GiftCardDTO is an structure to get api input in gift card api.
type CampaignDTO struct {
	ID    int                            `json:"id,string,omitempty"`
	Title string                         `json:"title"`
	Error *indraframework.IndraException `json:"error"`
}

func EmptyCampaignDTO() CampaignDTO {
	return CampaignDTO{}
}

func (a *CampaignDTO) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
