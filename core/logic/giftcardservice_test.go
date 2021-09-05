package logic_test

import (
	"errors"
	"giftcard-engine/core"
	"giftcard-engine/core/common"
	"giftcard-engine/core/dbmodel"
	"giftcard-engine/core/dto"
	"giftcard-engine/core/logic"
	"giftcard-engine/infrastructure/repository/sql"
	"giftcard-engine/utils/date"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

const (
	defaultBehavior = iota
	notFound
	internalError
	emptyData
	invalidOperation
)

type fakeGiftCardRepo struct {
	findByUUNCall       int32
	findByIDCall        int32
	storeCall           int32
	deleteCall          int32
	findByPublicKeyCall int32
	findPageCall        int32
	findBySecretKeyCall int32
	rollBackApproveCall int32
	strategy            int
}

var fakeInternalError = errors.New("repository internal error")

func (f *fakeGiftCardRepo) FindByUUN(uun string) []dbmodel.GiftCard {
	atomic.AddInt32(&f.findByUUNCall, 1)
	if f.strategy == emptyData {
		return []dbmodel.GiftCard{}
	}
	return []dbmodel.GiftCard{
		{Amount: int32(2000), PublicCode: "public",
			SecretCode: "secret", UUN: "milawd", ExpireDate: time.Now().UTC(), Status: dbmodel.Empty},
	}
}

func (f *fakeGiftCardRepo) FindByID(id uint) (*dbmodel.GiftCard, error) {
	atomic.AddInt32(&f.findByIDCall, 1)
	if f.strategy == notFound {
		return nil, common.GiftCardNotFound
	}
	return &dbmodel.GiftCard{
		Amount:     2000,
		PublicCode: "123456789012",
		SecretCode: "1234567890123456",
		UUN:        "",
		ExpireDate: date.DefaultToTimeOrDefault("2400-02-02"),
		Status:     dbmodel.Empty,
	}, nil
}

func (f *fakeGiftCardRepo) Store(card *dbmodel.GiftCard) error {
	atomic.AddInt32(&f.storeCall, 1)
	if f.strategy == internalError {
		return fakeInternalError
	}
	//NOTE: just for one specific scenario
	if card.Amount <= 0 {
		return fakeInternalError
	}
	return nil
}

func (f *fakeGiftCardRepo) Delete(card dbmodel.GiftCard) error {
	atomic.AddInt32(&f.deleteCall, 1)
	if f.strategy == notFound {
		return common.GiftCardNotFound
	}
	return nil
}

func (f *fakeGiftCardRepo) FindByPublicKey(key string) (*dbmodel.GiftCard, error) {
	atomic.AddInt32(&f.findByPublicKeyCall, 1)
	if f.strategy == notFound {
		return nil, common.GiftCardNotFound
	}
	return &dbmodel.GiftCard{
		Amount:     10,
		PublicCode: "public",
		SecretCode: "secret",
		UUN:        "milawd",
		ExpireDate: time.Now().UTC(),
		Status:     dbmodel.Approved,
	}, nil
}

func (f *fakeGiftCardRepo) FindPage(size, number uint, search string, campaignId *int,
	isValid *bool, expireDateFrom *time.Time, expireDateTo *time.Time) ([]dbmodel.GiftCard, int) {
	atomic.AddInt32(&f.findPageCall, 1)
	return []dbmodel.GiftCard{}, 0
}

func (f *fakeGiftCardRepo) FindBySecretKey(secret string) (*dbmodel.GiftCard, error) {
	atomic.AddInt32(&f.findBySecretKeyCall, 1)
	if f.strategy == notFound {
		return nil, common.GiftCardNotFound
	}
	return &dbmodel.GiftCard{
		Amount:     2000,
		PublicCode: "public",
		SecretCode: secret,
		UUN:        "",
		ExpireDate: time.Now().Add(25 * time.Hour),
		Status:     dbmodel.Empty,
	}, nil
}

func (f *fakeGiftCardRepo) RollBackApprove(secret string) error {
	atomic.AddInt32(&f.rollBackApproveCall, 1)

	if f.strategy == notFound {
		return common.GiftCardNotFound
	}

	if f.strategy == internalError {
		return fakeInternalError
	}
	return nil
}

func newFakeGiftCardRepo(strategy int) *fakeGiftCardRepo {
	return &fakeGiftCardRepo{
		strategy: strategy,
	}
}

/////////////////////////////////////
type fakeGiftCardMapper struct {
	ToGiftCardCall                  int32
	ToGiftCardDTOCall               int32
	ToListOfGiftCardDTOCall         int32
	ToGiftCardStatusDTOCall         int32
	ApprovedToGiftCardStatusDTOCall int32
	ToCampaignCall                  int32
	ToCampaignDTOCall               int32
	ToListOfCampaignsCall           int32
	actualMapper                    core.Mapper
}

func (f *fakeGiftCardMapper) ToGiftCard(dto dto.CreateGiftCardDTO) *dbmodel.GiftCard {
	atomic.AddInt32(&f.ToGiftCardCall, 1)
	return f.actualMapper.ToGiftCard(dto)
}
func (f *fakeGiftCardMapper) ToGiftCardDTO(card *dbmodel.GiftCard) dto.GiftCardDTO {
	atomic.AddInt32(&f.ToGiftCardDTOCall, 1)
	return f.actualMapper.ToGiftCardDTO(card)
}
func (f *fakeGiftCardMapper) ToListOfGiftCardDTO(cards []dbmodel.GiftCard) *dto.GiftCardsListDTO {
	atomic.AddInt32(&f.ToListOfGiftCardDTOCall, 1)
	return f.actualMapper.ToListOfGiftCardDTO(cards)
}
func (f *fakeGiftCardMapper) ToGiftCardStatusDTO(card dbmodel.GiftCard) dto.GiftCardStatusDTO {
	atomic.AddInt32(&f.ToGiftCardStatusDTOCall, 1)
	return f.actualMapper.ToGiftCardStatusDTO(card)
}
func (f *fakeGiftCardMapper) ApprovedToGiftCardStatusDTO(card dbmodel.GiftCard) dto.GiftCardStatusDTO {
	atomic.AddInt32(&f.ApprovedToGiftCardStatusDTOCall, 1)
	return f.actualMapper.ApprovedToGiftCardStatusDTO(card)
}

func (f *fakeGiftCardMapper) ToCampaign(dto dto.CreateCampaignDTO) dbmodel.Campaign {
	atomic.AddInt32(&f.ToCampaignCall, 1)
	return f.actualMapper.ToCampaign(dto)
}

func (f *fakeGiftCardMapper) ToCampaignDTO(campaign dbmodel.Campaign) dto.CampaignDTO {
	atomic.AddInt32(&f.ToCampaignDTOCall, 1)
	return f.actualMapper.ToCampaignDTO(campaign)
}

func (f *fakeGiftCardMapper) ToListOfCampaigns(campaigns []dbmodel.Campaign) []dto.CampaignDTO {
	atomic.AddInt32(&f.ToListOfCampaignsCall, 1)
	return f.actualMapper.ToListOfCampaigns(campaigns)
}

func newFakeGiftCardMapper() *fakeGiftCardMapper {
	return &fakeGiftCardMapper{
		actualMapper: sql.NewMapper(),
	}
}

//////end of fake dependencies

func createServiceForTest(strategy int) (core.GiftCardService, *fakeGiftCardRepo, *fakeGiftCardMapper) {
	mapper := newFakeGiftCardMapper()
	repo := newFakeGiftCardRepo(strategy)
	return logic.NewGiftCardService(repo, mapper), repo, mapper
}

func TestFindByUUN(te *testing.T) {
	te.Parallel()

	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		cards, err := service.FindByUUN("someone")

		assert.Empty(t, err)
		assert.Equal(t, int32(1), repo.findByUUNCall)
		assert.Equal(t, 1, len(cards.Cards))
		assert.Nil(t, cards.Error)
		assert.Equal(t, int32(1), mapper.ToListOfGiftCardDTOCall)
	})

	te.Run("with empty data strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(emptyData)

		cards, err := service.FindByUUN("someone")

		assert.NotEmpty(t, err)
		assert.Equal(t, common.NoGiftCardFoundForUser, err)
		assert.Equal(t, int32(1), repo.findByUUNCall)
		assert.Empty(t, cards)
		assert.Equal(t, int32(0), mapper.ToListOfGiftCardDTOCall)
	})
}

func TestFindByID(te *testing.T) {
	te.Parallel()
	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		card, err := service.FindByID(123)

		assert.Empty(t, err)
		assert.NotEmpty(t, card)
		assert.Equal(t, int32(1), repo.findByIDCall)
		assert.Equal(t, int32(1), mapper.ToGiftCardDTOCall)
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(notFound)

		card, err := service.FindByID(123)

		assert.NotEmpty(t, err)
		assert.Empty(t, card)
		assert.Equal(t, common.GiftCardNotFound, err)
		assert.Equal(t, int32(1), repo.findByIDCall)
		assert.Equal(t, int32(0), mapper.ToGiftCardDTOCall)
	})
}

func TestFindByPublicKey(te *testing.T) {
	te.Parallel()
	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		card, err := service.FindByPublicKey("123456123456")

		assert.Empty(t, err)
		assert.NotEmpty(t, card)
		assert.Equal(t, int32(1), repo.findByPublicKeyCall)
		assert.Equal(t, int32(1), mapper.ToGiftCardStatusDTOCall)
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(notFound)

		card, err := service.FindByPublicKey("123456123456")

		assert.NotEmpty(t, err)
		assert.Empty(t, card)
		assert.Equal(t, common.GiftCardNotFound, err)
		assert.Equal(t, int32(1), repo.findByPublicKeyCall)
		assert.Equal(t, int32(0), mapper.ToGiftCardStatusDTOCall)
	})
}

func TestStore(te *testing.T) {
	te.Parallel()
	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		card, err := service.Store(&dto.CreateGiftCardDTO{
			ExpireDate: "2300-02-02",
			Amount:     2000,
		})

		assert.Empty(t, err)
		assert.NotEmpty(t, card)
		assert.Equal(t, int32(1), repo.storeCall)
		assert.Equal(t, int32(1), mapper.ToGiftCardCall)
		assert.Equal(t, int32(1), mapper.ToGiftCardDTOCall)
	})

	te.Run("with internal server error strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(internalError)

		card, err := service.Store(&dto.CreateGiftCardDTO{
			ExpireDate: "2300-02-02",
			Amount:     2000,
		})

		assert.NotEmpty(t, err)
		assert.Empty(t, card)
		assert.Equal(t, int32(1), repo.storeCall)
		assert.Equal(t, int32(1), mapper.ToGiftCardCall)
		assert.Equal(t, int32(0), mapper.ToGiftCardDTOCall)
	})
}

func TestUpdate(te *testing.T) {
	te.Parallel()
	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		card, err := service.Update(&dto.UpdateGiftCardDto{
			ExpireDate: "2300-02-02",
			Amount:     2000,
			ID:         10,
		})

		assert.Empty(t, err)
		assert.NotEmpty(t, card)
		assert.Equal(t, int32(1), repo.storeCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
		assert.Equal(t, int32(1), mapper.ToGiftCardDTOCall)
	})

	te.Run("with invalid date string", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		card, err := service.Update(&dto.UpdateGiftCardDto{
			ExpireDate: "invalid date",
			Amount:     2000,
			ID:         10,
		})

		assert.NotEmpty(t, err)
		assert.Empty(t, card)
		assert.Equal(t, int32(0), repo.storeCall)
		assert.Equal(t, int32(0), repo.findByIDCall)
		assert.Equal(t, int32(0), mapper.ToGiftCardDTOCall)
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(notFound)

		card, err := service.Update(&dto.UpdateGiftCardDto{
			ExpireDate: "2400-02-02",
			Amount:     2000,
			ID:         10,
		})

		assert.NotEmpty(t, err)
		assert.Empty(t, card)
		assert.Equal(t, common.GiftCardNotFound, err)
		assert.Equal(t, int32(0), repo.storeCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
		assert.Equal(t, int32(0), mapper.ToGiftCardDTOCall)
	})

	te.Run("with internal server error strategy", func(t *testing.T) {
		service, repo, mapper := createServiceForTest(internalError)

		card, err := service.Update(&dto.UpdateGiftCardDto{
			ExpireDate: "2400-02-02",
			Amount:     2000,
			ID:         10,
		})

		assert.NotEmpty(t, err)
		assert.NotEmpty(t, card)
		assert.Equal(t, int32(1), repo.storeCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
		assert.Equal(t, int32(1), mapper.ToGiftCardDTOCall)
	})
}

func TestDelete(te *testing.T) {
	te.Parallel()
	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, _ := createServiceForTest(defaultBehavior)

		err := service.Delete(123)

		assert.Empty(t, err)
		assert.Equal(t, int32(1), repo.deleteCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, _ := createServiceForTest(notFound)

		err := service.Delete(123)

		assert.NotEmpty(t, err)
		assert.Equal(t, common.GiftCardNotFound, err)
		assert.Equal(t, int32(0), repo.deleteCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
	})
}

func TestCreateMany(te *testing.T) {
	te.Parallel()
	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)
		bulkDto := dto.BulkCreateGiftCardsDTO{GiftCards: []dto.CreateGiftCardDTO{
			{ExpireDate: "2400-02-02", Amount: 3000},
			{ExpireDate: "2400-02-10", Amount: 4000},
		}}

		cards, err := service.CreateMany(&bulkDto)

		assert.Empty(t, err)
		assert.NotEmpty(t, cards)
		assert.Equal(t, 2, len(cards.Cards))
		assert.Equal(t, int32(2), repo.storeCall)
		assert.Equal(t, int32(2), mapper.ToGiftCardDTOCall)
	})

	te.Run("with internal error strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(internalError)
		bulkDto := dto.BulkCreateGiftCardsDTO{GiftCards: []dto.CreateGiftCardDTO{
			{ExpireDate: "2400-02-02", Amount: 3000},
			{ExpireDate: "2400-02-10", Amount: 4000},
		}}

		cards, err := service.CreateMany(&bulkDto)

		assert.NotEmpty(t, err)
		assert.Empty(t, cards.Cards)
		assert.Equal(t, 0, len(cards.Cards))
		assert.Equal(t, int32(2), repo.storeCall)
		assert.Equal(t, int32(0), mapper.ToGiftCardDTOCall)
	})

	te.Run("try to add 2 object but one of the objects are valid", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)
		bulkDto := dto.BulkCreateGiftCardsDTO{GiftCards: []dto.CreateGiftCardDTO{
			{ExpireDate: "2400-02-02", Amount: -10},
			{ExpireDate: "2400-02-10", Amount: 4000},
		}}

		cards, err := service.CreateMany(&bulkDto)

		assert.NotEmpty(t, err)
		assert.NotEmpty(t, cards)
		assert.Equal(t, 1, len(cards.Cards))
		assert.Equal(t, int32(2), repo.storeCall)
		assert.Equal(t, int32(1), mapper.ToGiftCardDTOCall)
	})
}

func TestCreateSameMany(te *testing.T) {
	te.Parallel()
	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		cards, err := service.CreateSameMany(&dto.BulkCreateSameGiftCardsDTO{
			ExpireDate: "2400-02-02",
			Amount:     2000,
			Count:      20000,
		})

		assert.Empty(t, err)
		assert.NotEmpty(t, cards)
		assert.Equal(t, 20000, len(cards.Cards))
		assert.Equal(t, int32(20000), repo.storeCall)
		assert.Equal(t, int32(20000), mapper.ToGiftCardDTOCall)
	})

	te.Run("with internal error strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(internalError)

		cards, err := service.CreateSameMany(&dto.BulkCreateSameGiftCardsDTO{
			ExpireDate: "2400-02-02",
			Amount:     2000,
			Count:      200,
		})

		assert.NotEmpty(t, err)
		assert.Empty(t, cards.Cards)
		assert.Equal(t, 0, len(cards.Cards))
		assert.Equal(t, int32(200), repo.storeCall)
		assert.Equal(t, int32(0), mapper.ToGiftCardDTOCall)
	})
}

func TestFindPage(te *testing.T) {
	te.Parallel()
	service, repo, mapper := createServiceForTest(defaultBehavior)
	startDate := date.DefaultToTimeOrDefault("2050-01-01")
	endDate := date.DefaultToTimeOrDefault("2050-01-02")
	pageRes := service.FindPage(10, 10, "", nil, nil,
		&startDate, &endDate)

	assert.NotEmpty(te, pageRes)
	assert.Equal(te, 11, pageRes.Page)
	assert.Equal(te, 10, pageRes.Size)
	assert.Equal(te, 0, pageRes.TotalItems)
	assert.Empty(te, pageRes.GiftCards)
	assert.Equal(te, int32(1), mapper.ToListOfGiftCardDTOCall)
	assert.Equal(te, int32(1), repo.findPageCall)
}

func TestValidateGiftCards(te *testing.T) {
	te.Parallel()
	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		cards := service.ValidateGiftCards(&dto.ValidateGiftCardsDto{
			GiftCardsSecret: []string{
				"1234567890123456",
				"2234567890123456",
			},
		})

		assert.NotEmpty(t, cards)
		assert.Equal(t, 2, len(cards.Cards))
		assert.Equal(t, int32(2), repo.findBySecretKeyCall)
		assert.Equal(t, int32(2), mapper.ToGiftCardStatusDTOCall)
		assert.Equal(t, true, cards.Cards[0].IsValid)
		assert.Equal(t, true, cards.Cards[1].IsValid)
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(notFound)

		cards := service.ValidateGiftCards(&dto.ValidateGiftCardsDto{
			GiftCardsSecret: []string{
				"1234567890123456",
				"2234567890123456",
			},
		})

		assert.NotEmpty(t, cards)
		assert.Equal(t, 2, len(cards.Cards))
		assert.Equal(t, int32(2), repo.findBySecretKeyCall)
		assert.Equal(t, int32(0), mapper.ToGiftCardStatusDTOCall)
		assert.Equal(t, false, cards.Cards[0].IsValid)
		assert.Equal(t, false, cards.Cards[1].IsValid)
	})
}

func TestApproveGiftCards(te *testing.T) {
	te.Parallel()
	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		cards, err := service.ApproveGiftCards(&dto.ApproveGiftCardsDTO{
			UUN: "milawd",
			GiftCardsSecret: []string{
				"1234567890123456",
				"2234567890123456",
			},
		})

		assert.Empty(t, err)
		assert.NotEmpty(t, cards)
		assert.Equal(t, 2, len(cards.Cards))
		assert.Equal(t, int32(2), repo.findBySecretKeyCall)
		assert.Equal(t, int32(2), mapper.ApprovedToGiftCardStatusDTOCall)
		assert.Equal(t, true, cards.Cards[0].IsValid)
		assert.Equal(t, true, cards.Cards[1].IsValid)
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(notFound)

		cards, err := service.ApproveGiftCards(&dto.ApproveGiftCardsDTO{
			UUN: "milawd",
			GiftCardsSecret: []string{
				"1234567890123456",
				"2234567890123456",
			},
		})

		assert.NotEmpty(t, err)
		assert.Empty(t, cards)
		assert.Equal(t, err, common.GiftCardNotFound)
		assert.Equal(t, int32(2), repo.findBySecretKeyCall)
		assert.Equal(t, int32(0), mapper.ApprovedToGiftCardStatusDTOCall)
	})

	te.Run("with internal server error strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(internalError)

		cards, err := service.ApproveGiftCards(&dto.ApproveGiftCardsDTO{
			UUN: "milawd",
			GiftCardsSecret: []string{
				"1234567890123456",
				"2234567890123456",
			},
		})

		assert.NotEmpty(t, err)
		assert.Empty(t, cards)
		assert.Equal(t, int32(2), repo.findBySecretKeyCall)
		assert.Equal(t, int32(0), mapper.ApprovedToGiftCardStatusDTOCall)
	})
}

func TestValidateGiftCard(te *testing.T) {
	te.Parallel()
	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		card := service.ValidateGiftCard("1234567890123456")

		assert.NotEmpty(t, card)
		assert.Equal(t, int32(1), repo.findBySecretKeyCall)
		assert.Equal(t, int32(1), mapper.ToGiftCardStatusDTOCall)
		assert.Equal(t, true, card.IsValid)
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(notFound)

		card := service.ValidateGiftCard("1234567890123456")

		assert.NotEmpty(t, card)
		assert.Equal(t, int32(1), repo.findBySecretKeyCall)
		assert.Equal(t, int32(0), mapper.ToGiftCardStatusDTOCall)
		assert.Equal(t, false, card.IsValid)
	})
}

func TestApproveGiftCard(te *testing.T) {
	te.Parallel()

	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(defaultBehavior)

		card, err := service.ApproveGiftCard("milawd", "1234567890123456")

		assert.Empty(t, err)
		assert.NotEmpty(t, card)
		assert.Equal(t, int32(1), repo.findBySecretKeyCall)
		assert.Equal(t, int32(1), mapper.ApprovedToGiftCardStatusDTOCall)
		assert.Equal(t, true, card.IsValid)
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(notFound)

		card, err := service.ApproveGiftCard("milawd", "2234567890123456")

		assert.NotEmpty(t, err)
		assert.Empty(t, card)
		assert.Equal(t, err, common.GiftCardNotFound)
		assert.Equal(t, int32(1), repo.findBySecretKeyCall)
		assert.Equal(t, int32(0), mapper.ApprovedToGiftCardStatusDTOCall)
		assert.Equal(t, false, card.IsValid)
	})

	te.Run("with internal server error strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createServiceForTest(internalError)

		card, err := service.ApproveGiftCard("milawd", "2234567890123456")

		assert.NotEmpty(t, err)
		assert.Empty(t, card)
		assert.Equal(t, int32(1), repo.findBySecretKeyCall)
		assert.Equal(t, int32(0), mapper.ApprovedToGiftCardStatusDTOCall)
	})
}
