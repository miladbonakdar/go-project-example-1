package sql_test

import (
	"giftcard-engine/core/dbmodel"
	"giftcard-engine/core/dto"
	"giftcard-engine/infrastructure/repository/sql"
	"giftcard-engine/utils/date"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mapper = sql.NewMapper()

func TestToGiftCard(t *testing.T) {
	t.Parallel()
	cardDto := dto.CreateGiftCardDTO{
		ExpireDate: "2012-01-01",
		Amount:     2000,
	}

	card := mapper.ToGiftCard(cardDto)
	d, _ := date.DefaultToTime(cardDto.ExpireDate)

	assert.NotEmpty(t, card)
	assert.Equal(t, cardDto.Amount, card.Amount)
	assert.Equal(t, d, card.ExpireDate)
}

func TestToGiftCardDTO(t *testing.T) {
	t.Parallel()
	d, _ := date.DefaultToTime("2012-01-01")
	card := dbmodel.GiftCard{
		Amount:     2000,
		PublicCode: "public",
		SecretCode: "secret",
		UUN:        "milawd",
		ExpireDate: d,
		Status:     dbmodel.Empty,
		CampaignId: 1,
	}

	cardDto := mapper.ToGiftCardDTO(&card)

	assert.NotEmpty(t, cardDto)
	assert.Equal(t, card.Amount, cardDto.Amount)
	assert.Equal(t, card.SecretCode, cardDto.SecretCode)
	assert.Equal(t, card.PublicCode, cardDto.PublicCode)
	assert.Equal(t, card.CampaignId, cardDto.CampaignId)
	assert.Equal(t, card.UUN, cardDto.UUN)
	assert.Equal(t, card.ExpireDate.String(), "2012-01-01 00:00:00 +0000 UTC")
}

func TestToListOfGiftCardDTO(t *testing.T) {
	t.Parallel()
	d, _ := date.DefaultToTime("2012-01-01")
	cards := []dbmodel.GiftCard{
		{
			Amount:     2000,
			PublicCode: "public",
			SecretCode: "secret",
			UUN:        "milawd",
			ExpireDate: d,
			Status:     dbmodel.Empty,
			CampaignId: 1,
		},
	}

	cardsDto := mapper.ToListOfGiftCardDTO(cards)

	assert.Equal(t, len(cards), len(cardsDto.Cards))
	assert.NotEmpty(t, cardsDto)
	assert.Equal(t, cards[0].Amount, cardsDto.Cards[0].Amount)
	assert.Equal(t, cards[0].SecretCode, cardsDto.Cards[0].SecretCode)
	assert.Equal(t, cards[0].PublicCode, cardsDto.Cards[0].PublicCode)
	assert.Equal(t, cards[0].CampaignId, cardsDto.Cards[0].CampaignId)
	assert.Equal(t, cards[0].UUN, cardsDto.Cards[0].UUN)
	assert.Equal(t, cards[0].ExpireDate.String(), "2012-01-01 00:00:00 +0000 UTC")
}

func TestToGiftCardStatusDTO(t *testing.T) {
	t.Parallel()
	d, _ := date.DefaultToTime("2222-01-01")
	card := dbmodel.GiftCard{
		Amount:     2000,
		PublicCode: "public",
		SecretCode: "secret",
		UUN:        "",
		ExpireDate: d,
		Status:     dbmodel.Empty,
		CampaignId: 1,
	}

	statusDto := mapper.ToGiftCardStatusDTO(card)

	assert.NotEmpty(t, statusDto)
	assert.Equal(t, card.Amount, statusDto.Amount)
	assert.Equal(t, true, statusDto.IsValid)
	assert.Equal(t, card.SecretCode, statusDto.SecretKey)
}

func TestApprovedToGiftCardStatusDTO(t *testing.T) {
	t.Parallel()
	d, _ := date.DefaultToTime("2222-01-01")
	card := dbmodel.GiftCard{
		Amount:     2000,
		PublicCode: "public",
		SecretCode: "secret",
		UUN:        "",
		ExpireDate: d,
		Status:     dbmodel.Empty,
		CampaignId: 1,
	}

	statusDto := mapper.ApprovedToGiftCardStatusDTO(card)

	assert.NotEmpty(t, statusDto)
	assert.Equal(t, card.Amount, statusDto.Amount)
	assert.Equal(t, true, statusDto.IsValid)
	assert.Equal(t, card.SecretCode, statusDto.SecretKey)
}

func TestToCampaign(t *testing.T) {
	t.Parallel()
	camp := dto.CreateCampaignDTO{
		Title: "dastan",
	}

	model := mapper.ToCampaign(camp)

	assert.NotEmpty(t, model)
	assert.Equal(t, camp.Title, model.Title)
}

func TestToCampaignDTO(t *testing.T) {
	t.Parallel()
	camp := dbmodel.Campaign{
		Title: "",
	}
	camp.ID = 1

	Dto := mapper.ToCampaignDTO(camp)

	assert.NotEmpty(t, Dto)
	assert.Equal(t, camp.Title, Dto.Title)
	assert.Equal(t, camp.ID, Dto.ID)
}

func TestToListOfCampaigns(t *testing.T) {
	t.Parallel()
	camps := []dbmodel.Campaign{
		{
			Title: "dastan",
		},
	}
	camps[0].ID = 1

	campsDto := mapper.ToListOfCampaigns(camps)

	assert.Equal(t, len(camps), len(campsDto))
	assert.NotEmpty(t, campsDto)
	assert.Equal(t, camps[0].ID, campsDto[0].ID)
	assert.Equal(t, camps[0].Title, campsDto[0].Title)
}
