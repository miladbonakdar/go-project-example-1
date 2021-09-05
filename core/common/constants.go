package common

import (
	"errors"
)

var (
	GiftCardNotFound            = errors.New("gift card cannot be found")
	CampaignNotFound            = errors.New("campaign cannot be found")
	NoGiftCardFoundForUser      = errors.New("there are no gift cards for this user")
	GiftCardIsTaken             = errors.New("the Gift card is taken by another user")
	GiftCardIsNotValid          = errors.New("this gift card is not valid anymore. you cannot update it")
	ExpireDateIsNotInValidRange = errors.New("the expire date should be after today")
	DuplicatedCampaignTitle     = errors.New("duplicated campaign name")
	InvalidCampaign             = errors.New("invalid campaign")
	InvalidCampaignQueryParam   = errors.New("invalid campaign query param")
	InvalidIsValidQueryParam    = errors.New("invalid isValid query param")
)
