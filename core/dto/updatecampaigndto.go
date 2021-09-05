package dto

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateCampaignDto struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (a UpdateCampaignDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title, validation.Required),
		validation.Field(&a.ID, validation.Required),
	)
}
