package logic

import (
	"giftcard-engine/core"
	"giftcard-engine/core/dto"
	"giftcard-engine/infrastructure/logger"
)

type campaignService struct {
	repo   core.CampaignRepository
	mapper core.Mapper
}

func (g *campaignService) Create(campaign dto.CreateCampaignDTO) (dto.CampaignDTO, error) {
	c := g.mapper.ToCampaign(campaign)
	err := g.repo.Store(&c)
	campaignDto := g.mapper.ToCampaignDTO(c)
	if err != nil {
		logger.ErrorException(err, "error in creating new campaign")
		return dto.EmptyCampaignDTO(), err
	}
	return campaignDto, nil
}

func (g *campaignService) Update(campaign dto.UpdateCampaignDto) (dto.CampaignDTO, error) {
	campaignModel, err := g.repo.FindByID(uint(campaign.ID))
	if err != nil {
		logger.WithData(campaign).ErrorException(err, "error in updating a campaign")
		return dto.EmptyCampaignDTO(), err
	}
	(&campaignModel).Update(campaign.Title)
	campaignDto := g.mapper.ToCampaignDTO(campaignModel)
	return campaignDto, g.repo.Store(&campaignModel)
}

func (g *campaignService) Delete(id uint) error {
	campaign, err := g.repo.FindByID(id)
	if err != nil {
		logger.WithData(map[string]interface{}{
			"id" : id,
		}).ErrorException(err, "error in deleting a campaign")
		return err
	}
	return g.repo.Delete(campaign)
}

func (g *campaignService) FindPage(size, page uint, search string) dto.CampaignPageDTO {
	campaigns, total := g.repo.FindPage(size, page, search)
	return dto.NewCampaignPageDTO(g.mapper.ToListOfCampaigns(campaigns), int(size), int(page), total)
}

func NewCampaignService(repository core.CampaignRepository, mapper core.Mapper) core.CampaignService {
	return &campaignService{repo: repository, mapper: mapper}
}
