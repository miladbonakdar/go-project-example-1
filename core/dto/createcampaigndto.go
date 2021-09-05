package dto

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type CreateCampaignDTO struct {
	Title string `json:"title"`
}

func (a CreateCampaignDTO) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title, validation.Required),
	)
}
