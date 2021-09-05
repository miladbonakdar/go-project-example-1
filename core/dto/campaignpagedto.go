package dto

import "giftcard-engine/utils/indraframework"

type CampaignPageDTO struct {
	Size       int                            `json:"size"`
	Page       int                            `json:"page"`
	Campaigns  []CampaignDTO                  `json:"campaigns"`
	TotalItems int                            `json:"total_items"`
	Error      *indraframework.IndraException `json:"error"`
}

func NewCampaignPageDTO(campaigns []CampaignDTO, size, page, total int) CampaignPageDTO {
	return CampaignPageDTO{
		Size:       size,
		Page:       page + 1,
		Campaigns:  campaigns,
		TotalItems: total,
	}
}

func (a *CampaignPageDTO) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
