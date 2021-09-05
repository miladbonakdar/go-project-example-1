package logic

import (
	"giftcard-engine/core"
	"giftcard-engine/core/common"
	"giftcard-engine/core/dbmodel"
	"giftcard-engine/core/dto"
	"giftcard-engine/infrastructure/logger"
	"giftcard-engine/utils/date"
	"strings"
	"sync"
	"time"
)

type giftCardService struct {
	giftCardRepo core.GiftCardRepository
	mapper       core.Mapper
}

func (g *giftCardService) FindByUUN(uun string) (*dto.GiftCardsListDTO, error) {
	giftCards := g.giftCardRepo.FindByUUN(uun)
	if len(giftCards) == 0 {
		return nil, common.NoGiftCardFoundForUser
	}
	return g.mapper.ToListOfGiftCardDTO(giftCards), nil
}

// FindByID find a gift card by ID
func (g *giftCardService) FindByID(id uint) (*dto.GiftCardDTO, error) {
	giftCard, err := g.giftCardRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	giftCardDto := g.mapper.ToGiftCardDTO(giftCard)
	return &giftCardDto, nil
}

func (g *giftCardService) FindByPublicKey(key string) (*dto.GiftCardStatusDTO, error) {
	giftCard, err := g.giftCardRepo.FindByPublicKey(key)
	if err != nil {
		return nil, err
	}
	giftCardStatusDto := g.mapper.ToGiftCardStatusDTO(*giftCard)
	return &giftCardStatusDto, nil
}

// Store store a gift card
func (g *giftCardService) Store(card *dto.CreateGiftCardDTO) (*dto.GiftCardDTO, error) {
	giftCard := g.mapper.ToGiftCard(*card)
	err := g.giftCardRepo.Store(giftCard)
	if err != nil {
		logger.WithData(card).ErrorException(err,"error while storing a gift card")
		return nil, err
	}
	giftCardDto := g.mapper.ToGiftCardDTO(giftCard)
	return &giftCardDto, nil
}

// Store store a gift card
func (g *giftCardService) Update(card *dto.UpdateGiftCardDto) (*dto.GiftCardDTO, error) {
	expDate, err := date.DefaultToTime(card.ExpireDate)
	if err != nil {
		return nil, err
	}
	giftCard, err := g.giftCardRepo.FindByID(uint(card.ID))
	if err != nil {
		return nil, err
	}
	err = giftCard.Update(card.Amount, expDate)
	if err != nil {
		logger.WithData(card).ErrorException(err,"error while updating a gift card")
		return nil, err
	}
	giftCardDto := g.mapper.ToGiftCardDTO(giftCard)
	return &giftCardDto, g.giftCardRepo.Store(giftCard)
}

func (g *giftCardService) Delete(id uint) error {
	card, err := g.giftCardRepo.FindByID(id)
	if err != nil {
		return err
	}
	return g.giftCardRepo.Delete(*card)
}

func (g *giftCardService) CreateMany(cards *dto.BulkCreateGiftCardsDTO) (*dto.GiftCardsListDTO, error) {
	c := make(chan dto.GiftCardDTO, len(cards.GiftCards))
	errorChannel := make(chan error, len(cards.GiftCards))
	defer close(c)
	defer close(errorChannel)
	cardsLength := len(cards.GiftCards)
	for i := 0; i < cardsLength; i++ {
		card := cards.GiftCards[i]
		go g.createGiftCard(card.ExpireDate, card.Amount, card.CampaignId, c, errorChannel)
	}

	cardsDto := make([]dto.GiftCardDTO, 0, cardsLength)

	var err error
	for i := 0; i < cardsLength; i++ {
		select {
		case item := <-c:
			cardsDto = append(cardsDto, item)
		case err = <-errorChannel:
		}
	}

	return &dto.GiftCardsListDTO{
		Cards: cardsDto,
		Error: nil,
	}, err
}

func (g *giftCardService) CreateSameMany(cards *dto.BulkCreateSameGiftCardsDTO) (*dto.GiftCardsListDTO, error) {
	c := make(chan dto.GiftCardDTO, cards.Count)
	errorChannel := make(chan error, cards.Count)
	defer close(c)
	defer close(errorChannel)

	for i := 0; i < cards.Count; i++ {
		go g.createGiftCard(cards.ExpireDate, cards.Amount, cards.CampaignId, c, errorChannel)
	}

	cardsDto := make([]dto.GiftCardDTO, 0, cards.Count)

	var err error
	for i := 0; i < cards.Count; i++ {
		select {
		case item := <-c:
			cardsDto = append(cardsDto, item)
		case err = <-errorChannel:
		}
	}

	return &dto.GiftCardsListDTO{
		Cards: cardsDto,
		Error: nil,
	}, err
}

func (g *giftCardService) FindPage(size, page uint, search string, campaignId *int, isValid *bool,
	expireDateFrom *time.Time, expireDateTo *time.Time) dto.GiftCardsPageDTO {
	cards, total := g.giftCardRepo.FindPage(size, page, search, campaignId, isValid, expireDateFrom, expireDateTo)
	return *dto.NewGiftCardsPageDTO(g.mapper.ToListOfGiftCardDTO(cards).Cards, int(size), int(page), total)
}

func (g *giftCardService) ValidateGiftCards(validateDto *dto.ValidateGiftCardsDto) *dto.GiftCardStatusListDTO {
	secretsCount := len(validateDto.GiftCardsSecret)
	c := make(chan dto.GiftCardStatusDTO, secretsCount)
	defer close(c)

	for i := 0; i < secretsCount; i++ {
		go g.validateSecretKey(validateDto.GiftCardsSecret[i], c)
	}

	cardsDto := make([]dto.GiftCardStatusDTO, 0, secretsCount)

	for i := 0; i < secretsCount; i++ {
		cardsDto = append(cardsDto, <-c)
	}

	return &dto.GiftCardStatusListDTO{
		Cards: cardsDto,
		Error: nil,
	}
}

func (g *giftCardService) ApproveGiftCards(approveDto *dto.ApproveGiftCardsDTO) (*dto.GiftCardStatusListDTO, error) {
	secretsCount := len(approveDto.GiftCardsSecret)
	c := make(chan dto.GiftCardStatusDTO, secretsCount)
	errorChannel := make(chan error)
	defer close(c)
	defer close(errorChannel)

	for i := 0; i < secretsCount; i++ {
		go g.approveUser(approveDto.UUN, approveDto.GiftCardsSecret[i], c, errorChannel)
	}

	doneSecrets := make([]dto.GiftCardStatusDTO, 0, secretsCount)
	var err error
	for i := 0; i < secretsCount; i++ {
		select {
		case secret := <-c:
			doneSecrets = append(doneSecrets, secret)
		case err = <-errorChannel:
		}
	}
	if err != nil {
		g.rollBackApprovedCards(doneSecrets)
		return nil, err
	}

	return &dto.GiftCardStatusListDTO{
		Cards: doneSecrets,
		Error: nil,
	}, nil
}

func (g *giftCardService) ValidateGiftCard(giftCardSecret string) dto.GiftCardStatusDTO {
	c := make(chan dto.GiftCardStatusDTO)
	defer close(c)
	go g.validateSecretKey(giftCardSecret, c)
	return <-c
}

func (g *giftCardService) ApproveGiftCard(uun, giftCardSecret string) (dto.GiftCardStatusDTO, error) {
	c := make(chan dto.GiftCardStatusDTO)
	errorChannel := make(chan error)
	defer close(c)
	defer close(errorChannel)

	go g.approveUser(uun, giftCardSecret, c, errorChannel)

	select {
	case secret := <-c:
		return secret, nil
	case err := <-errorChannel:
		return dto.GiftCardStatusDTO{}, err
	}
}

func (g *giftCardService) rollBackApprovedCards(doneSecrets []dto.GiftCardStatusDTO) {
	wg := &sync.WaitGroup{}
	wg.Add(len(doneSecrets))
	for _, card := range doneSecrets {
		go func() {
			err := g.giftCardRepo.RollBackApprove(card.SecretKey)
			if err != nil {
				logger.ErrorException(err,"error while rolling back a gift card")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func (g *giftCardService) createGiftCard(expireDate string, amount int32, campaignId uint, channel chan<- dto.GiftCardDTO,
	errorChannel chan<- error) {
	giftCard := dbmodel.NewGiftCard(amount, date.DefaultToTimeOrDefault(expireDate))
	giftCard.SetCampaign(campaignId)
	for {
		err := g.giftCardRepo.Store(giftCard)
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "duplicate") {
			logger.ErrorException(err,"error while creating a new gift card")
			errorChannel <- err
			return
		}
		giftCard.GenerateKey()
	}
	channel <- g.mapper.ToGiftCardDTO(giftCard)
}

func (g *giftCardService) validateSecretKey(secret string, c chan<- dto.GiftCardStatusDTO) {
	secret = strings.ToUpper(secret)
	card, err := g.giftCardRepo.FindBySecretKey(secret)
	if err != nil {
		c <- dto.GiftCardStatusDTO{
			IsValid:   false,
			SecretKey: secret,
		}
		return
	}
	c <- g.mapper.ToGiftCardStatusDTO(*card)
}

func (g *giftCardService) approveUser(uun, secret string, data chan<- dto.GiftCardStatusDTO, errorChannel chan<- error) {
	secret = strings.ToUpper(secret)
	card, err := g.giftCardRepo.FindBySecretKey(secret)
	if err != nil {
		errorChannel <- err
		return
	}
	err = card.SetUUN(uun)
	if err != nil {
		logger.Error(err.Error())
		errorChannel <- err
		return
	}
	err = g.giftCardRepo.Store(card)
	if err != nil {
		errorChannel <- err
		logger.ErrorException(err,"error while approving a gift card")
		return
	}

	data <- g.mapper.ApprovedToGiftCardStatusDTO(*card)
}

func NewGiftCardService(repository core.GiftCardRepository, mapper core.Mapper) core.GiftCardService {
	return &giftCardService{giftCardRepo: repository, mapper: mapper}
}
