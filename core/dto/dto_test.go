package dto_test

import (
	"giftcard-engine/core/common"
	"giftcard-engine/core/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateApproveGiftCardsDTO(te *testing.T) {
	te.Parallel()
	te.Run("valid approveGiftCardsDTO", func(t *testing.T) {
		approveGiftCardsDTO := dto.ApproveGiftCardsDTO{
			UUN:             "milawd",
			GiftCardsSecret: []string{"123nlon123njkl12", "1234567890123456"},
		}
		err := approveGiftCardsDTO.Validate()
		assert.Empty(t, err)
	})

	te.Run("invalid secret key in the approveGiftCardsDTO", func(t *testing.T) {
		approveGiftCardsDTO := dto.ApproveGiftCardsDTO{
			UUN:             "milawd",
			GiftCardsSecret: []string{"invalid", "1234567890123456"},
		}
		err := approveGiftCardsDTO.Validate()
		assert.NotEmpty(t, err)
	})

	te.Run("empty uun in the approveGiftCardsDTO", func(t *testing.T) {
		approveGiftCardsDTO := dto.ApproveGiftCardsDTO{
			UUN:             "",
			GiftCardsSecret: []string{"1234567890123456"},
		}
		err := approveGiftCardsDTO.Validate()
		assert.NotEmpty(t, err)
	})
}

func TestValidateBulkCreateGiftCardsDTO(te *testing.T) {
	te.Parallel()
	item := dto.BulkCreateGiftCardsDTO{
		GiftCards: []dto.CreateGiftCardDTO{
			{
				ExpireDate: "2300-02-02",
				Amount:     2000,
				CampaignId: 1,
			},
		},
	}
	item2 := dto.BulkCreateGiftCardsDTO{
		GiftCards: []dto.CreateGiftCardDTO{
			{
				ExpireDate: "2300-02-02",
				Amount:     500,
				CampaignId: 1,
			},
		},
	}
	item3 := dto.BulkCreateGiftCardsDTO{
		GiftCards: []dto.CreateGiftCardDTO{
			{
				ExpireDate: "2300-02-02",
				Amount:     5000,
			},
		},
	}
	err := item.Validate()
	err2 := item2.Validate()
	err3 := item3.Validate()
	assert.Empty(te, err)
	assert.NotEmpty(te, err2)
	assert.NotEmpty(te, err3)
}

func TestValidateBulkCreateSameGiftCardsDTO(te *testing.T) {
	te.Parallel()
	te.Run("valid BulkCreateSameGiftCardsDTO", func(t *testing.T) {
		item := dto.BulkCreateSameGiftCardsDTO{
			ExpireDate: "2300-02-02",
			Amount:     10000,
			Count:      1000,
			CampaignId: 1,
		}
		err := item.Validate()
		assert.Empty(t, err)
	})

	te.Run("invalid amount in the BulkCreateSameGiftCardsDTO", func(t *testing.T) {
		item := dto.BulkCreateSameGiftCardsDTO{
			ExpireDate: "2300-02-02",
			Amount:     100,
			Count:      1000,
			CampaignId: 1,
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})

	te.Run("invalid campaign id BulkCreateSameGiftCardsDTO", func(t *testing.T) {
		item := dto.BulkCreateSameGiftCardsDTO{
			ExpireDate: "2300-02-02",
			Amount:     10000,
			Count:      1000,
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})

	te.Run("invalid expire date in the BulkCreateSameGiftCardsDTO", func(t *testing.T) {
		item := dto.BulkCreateSameGiftCardsDTO{
			ExpireDate: "1980-02-02",
			Amount:     10000,
			Count:      1000,
			CampaignId: 1,
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})

	te.Run("invalid count in the BulkCreateSameGiftCardsDTO", func(t *testing.T) {
		item := dto.BulkCreateSameGiftCardsDTO{
			ExpireDate: "2300-02-02",
			Amount:     10000,
			Count:      -10,
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})
}

func TestValidateCreateGiftCardDTO(te *testing.T) {
	te.Parallel()
	te.Run("valid CreateGiftCardDTO", func(t *testing.T) {
		item := dto.CreateGiftCardDTO{
			ExpireDate: "2300-02-02",
			Amount:     1000,
			CampaignId: 1,
		}
		err := item.Validate()
		assert.Empty(t, err)
	})

	te.Run("invalid expire date in CreateGiftCardDTO", func(t *testing.T) {
		item := dto.CreateGiftCardDTO{
			ExpireDate: "1980-02-02",
			Amount:     1000,
			CampaignId: 1,
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})

	te.Run("invalid amount in CreateGiftCardDTO", func(t *testing.T) {
		item := dto.CreateGiftCardDTO{
			ExpireDate: "2300-02-02",
			Amount:     10,
			CampaignId: 1,
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})
}

func TestValidateUpdateGiftCardDto(te *testing.T) {
	te.Parallel()
	te.Run("valid UpdateGiftCardDto", func(t *testing.T) {
		item := dto.UpdateGiftCardDto{
			ExpireDate: "2300-02-02",
			Amount:     1000,
			ID:         1000,
		}
		err := item.Validate()
		assert.Empty(t, err)
	})

	te.Run("invalid expire date in UpdateGiftCardDto", func(t *testing.T) {
		item := dto.UpdateGiftCardDto{
			ExpireDate: "1980-02-02",
			Amount:     1000,
			ID:         1000,
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})

	te.Run("invalid amount in UpdateGiftCardDto", func(t *testing.T) {
		item := dto.UpdateGiftCardDto{
			ExpireDate: "2300-02-02",
			Amount:     10,
			ID:         1000,
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})

	te.Run("invalid id in UpdateGiftCardDto", func(t *testing.T) {
		item := dto.UpdateGiftCardDto{
			ExpireDate: "2300-02-02",
			Amount:     10000,
			ID:         0,
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})
}

func TestValidateValidateGiftCardsDto(te *testing.T) {
	te.Parallel()
	te.Run("valid ValidateGiftCardsDto", func(t *testing.T) {
		item := dto.ValidateGiftCardsDto{
			GiftCardsSecret: []string{
				"1234567890123456",
				"sdfghjklkjhgfd12",
			},
		}
		err := item.Validate()
		assert.Empty(t, err)
	})

	te.Run("invalid secret key in the ValidateGiftCardsDto", func(t *testing.T) {
		item := dto.ValidateGiftCardsDto{
			GiftCardsSecret: []string{
				"1234567890123456",
				"invalid",
			},
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})
}

func TestCheckForDate(te *testing.T) {
	te.Parallel()
	te.Run("Valid date", func(t *testing.T) {
		err := dto.CheckForDate("2222-02-22")
		assert.Empty(t, err)
	})

	te.Run("inValid date that is before or equal to today date", func(t *testing.T) {
		err := dto.CheckForDate("1980-02-22")
		assert.NotEmpty(t, err)
		assert.Equal(t, common.ExpireDateIsNotInValidRange, err)
	})

	te.Run("inValid date string", func(t *testing.T) {
		err := dto.CheckForDate("invalid")
		assert.NotEmpty(t, err)
	})
}

func TestValidateCreateCampaignDTO(te *testing.T) {
	te.Parallel()
	te.Run("valid CreateCampaignDTO", func(t *testing.T) {
		item := dto.CreateCampaignDTO{
			Title: "dastan",
		}
		err := item.Validate()
		assert.Empty(t, err)
	})

	te.Run("invalid title in CreateCampaignDTO", func(t *testing.T) {
		item := dto.CreateCampaignDTO{}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})
}

func TestValidateUpdateCampaignDto(te *testing.T) {
	te.Parallel()
	te.Run("valid UpdateCampaignDto", func(t *testing.T) {
		item := dto.UpdateCampaignDto{
			Title: "dastan",
			ID:    1,
		}
		err := item.Validate()
		assert.Empty(t, err)
	})

	te.Run("invalid title in CreateCampaignDTO", func(t *testing.T) {
		item := dto.UpdateCampaignDto{
			ID: 1,
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})

	te.Run("invalid id in CreateCampaignDTO", func(t *testing.T) {
		item := dto.UpdateCampaignDto{
			Title: "dastan",
		}
		err := item.Validate()
		assert.NotEmpty(t, err)
	})
}
