package sql

import (
	"giftcard-engine/core"
	"giftcard-engine/core/dbmodel"
	"giftcard-engine/core/dto"
	"giftcard-engine/utils/date"
)

type mapper struct{}

func (m *mapper) ToGiftCard(dto dto.CreateGiftCardDTO) *dbmodel.GiftCard {
	t, _ := date.DefaultToTime(dto.ExpireDate)
	c := dbmodel.NewGiftCard(dto.Amount, t)
	c.SetCampaign(dto.CampaignId)
	return c
}

func (m *mapper) ToGiftCardDTO(card *dbmodel.GiftCard) dto.GiftCardDTO {

	if card.Campaign == nil {
		c := dbmodel.EmptyCampaign()
		card.Campaign = &c
	}
	return dto.GiftCardDTO{
		ID:            card.ID,
		PublicCode:    card.PublicCode,
		SecretCode:    card.SecretCode,
		UUN:           card.UUN,
		ExpireDate:    card.ExpireDate.Local().String(),
		Amount:        card.Amount,
		IsValid:       card.IsValid(),
		CampaignId:    card.CampaignId,
		CampaignTitle: card.Campaign.Title,
	}
}

func (m *mapper) ToListOfGiftCardDTO(cards []dbmodel.GiftCard) *dto.GiftCardsListDTO {
	var newList []dto.GiftCardDTO
	for _, card := range cards {
		newList = append(newList, m.ToGiftCardDTO(&card))
	}
	return &dto.GiftCardsListDTO{
		Cards: newList,
		Error: nil,
	}
}

func (m *mapper) ToGiftCardStatusDTO(card dbmodel.GiftCard) dto.GiftCardStatusDTO {
	return dto.GiftCardStatusDTO{
		Id:         card.ID,
		IsValid:    card.IsValid(),
		Amount:     card.Amount,
		SecretKey:  card.SecretCode,
		PublicKey:  card.PublicCode,
		UUN:        card.UUN,
		ExpireDate: card.ExpireDate.Local().String(),
	}
}

func (m *mapper) ApprovedToGiftCardStatusDTO(card dbmodel.GiftCard) dto.GiftCardStatusDTO {
	return dto.GiftCardStatusDTO{
		Id:         card.ID,
		IsValid:    card.IsDateValid(),
		SecretKey:  card.SecretCode,
		PublicKey:  card.PublicCode,
		Amount:     card.Amount,
		UUN:        card.UUN,
		ExpireDate: card.ExpireDate.Local().String(),
	}
}

func (m *mapper) ToCampaign(dto dto.CreateCampaignDTO) dbmodel.Campaign {
	return dbmodel.Campaign{
		Title: dto.Title,
	}
}

func (m *mapper) ToCampaignDTO(campaign dbmodel.Campaign) dto.CampaignDTO {
	return dto.CampaignDTO{
		ID:    campaign.ID,
		Title: campaign.Title,
		Error: nil,
	}
}

func (m *mapper) ToListOfCampaigns(campaigns []dbmodel.Campaign) []dto.CampaignDTO {
	var newList []dto.CampaignDTO
	for _, campaign := range campaigns {
		newList = append(newList, m.ToCampaignDTO(campaign))
	}
	return newList
}

func NewMapper() core.Mapper {
	return &mapper{}
}
