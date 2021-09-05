package logic_test

import (
	"giftcard-engine/core"
	"giftcard-engine/core/common"
	"giftcard-engine/core/dbmodel"
	"giftcard-engine/core/dto"
	"giftcard-engine/core/logic"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

type fakeCampaignRepo struct {
	findByIDCall int32
	storeCall    int32
	deleteCall   int32
	findPageCall int32
	strategy     int
}

var defaultCampaign = dbmodel.Campaign{
	Title: "dastan",
}

func (r *fakeCampaignRepo) FindByID(id uint) (dbmodel.Campaign, error) {
	atomic.AddInt32(&r.findByIDCall, 1)

	if r.strategy == notFound {
		return dbmodel.Campaign{}, common.CampaignNotFound
	}
	return defaultCampaign, nil
}

func (r *fakeCampaignRepo) Store(campaign *dbmodel.Campaign) error {
	atomic.AddInt32(&r.storeCall, 1)
	if r.strategy == internalError {
		return fakeInternalError
	}
	if r.strategy == invalidOperation {
		return common.DuplicatedCampaignTitle
	}
	return nil
}

func (r *fakeCampaignRepo) Delete(campaign dbmodel.Campaign) error {
	atomic.AddInt32(&r.deleteCall, 1)

	if r.strategy == notFound {
		return common.CampaignNotFound
	}

	if r.strategy == internalError {
		return fakeInternalError
	}

	return nil
}

func (r *fakeCampaignRepo) FindPage(size, number uint, search string) ([]dbmodel.Campaign, int) {
	atomic.AddInt32(&r.findPageCall, 1)
	return []dbmodel.Campaign{
		defaultCampaign,
	}, 1
}

func newFakeCampaignRepo(strategy int) *fakeCampaignRepo {
	return &fakeCampaignRepo{
		strategy: strategy,
	}
}

//////////////////

func createCampaignServiceForTest(strategy int) (core.CampaignService, *fakeCampaignRepo, *fakeGiftCardMapper) {
	mapper := newFakeGiftCardMapper()
	repo := newFakeCampaignRepo(strategy)
	return logic.NewCampaignService(repo, mapper), repo, mapper
}

func TestCampaignCreate(te *testing.T) {
	te.Parallel()

	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createCampaignServiceForTest(defaultBehavior)

		camp, err := service.Create(dto.CreateCampaignDTO{Title: "dastan"})

		assert.Empty(t, err)
		assert.Equal(t, int32(1), repo.storeCall)
		assert.NotNil(t, camp)
		assert.Equal(t, int32(1), mapper.ToCampaignCall)
	})

	te.Run("with duplicate title", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createCampaignServiceForTest(invalidOperation)

		_, err := service.Create(dto.CreateCampaignDTO{Title: "dastan"})

		assert.NotEmpty(t, err)
		assert.Equal(t, common.DuplicatedCampaignTitle, err)
		assert.Equal(t, int32(1), repo.storeCall)
		assert.Equal(t, int32(1), mapper.ToCampaignCall)
	})

	te.Run("with internal error", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createCampaignServiceForTest(internalError)

		_, err := service.Create(dto.CreateCampaignDTO{Title: "dastan"})

		assert.NotEmpty(t, err)
		assert.NotNil(t, err)
		assert.Equal(t, int32(1), repo.storeCall)
		assert.Equal(t, int32(1), mapper.ToCampaignCall)
	})
}

func TestCampaignUpdate(te *testing.T) {
	te.Parallel()

	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createCampaignServiceForTest(defaultBehavior)

		camp, err := service.Update(dto.UpdateCampaignDto{Title: "dastan",
			ID: 1})

		assert.Empty(t, err)
		assert.Equal(t, int32(1), repo.storeCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
		assert.NotNil(t, camp)
		assert.Equal(t, int32(1), mapper.ToCampaignDTOCall)
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createCampaignServiceForTest(notFound)

		_, err := service.Update(dto.UpdateCampaignDto{Title: "dastan",
			ID: 1})

		assert.NotNil(t, err)
		assert.Equal(t, int32(0), repo.storeCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
		assert.Equal(t, int32(0), mapper.ToCampaignDTOCall)
		assert.Equal(t, common.CampaignNotFound, err)
	})

	te.Run("with internal error strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, mapper := createCampaignServiceForTest(internalError)

		_, err := service.Update(dto.UpdateCampaignDto{Title: "dastan",
			ID: 1})

		assert.NotNil(t, err)
		assert.Equal(t, int32(1), repo.storeCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
		assert.Equal(t, int32(1), mapper.ToCampaignDTOCall)
	})
}

func TestCampaignDelete(te *testing.T) {
	te.Parallel()

	te.Run("default behavior", func(t *testing.T) {
		t.Parallel()
		service, repo, _ := createCampaignServiceForTest(defaultBehavior)

		err := service.Delete(12)

		assert.Empty(t, err)
		assert.Equal(t, int32(1), repo.deleteCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, _ := createCampaignServiceForTest(notFound)

		err := service.Delete(12)

		assert.NotNil(t, err)
		assert.Equal(t, int32(0), repo.deleteCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
		assert.Equal(t, common.CampaignNotFound, err)
	})

	te.Run("with internal error strategy", func(t *testing.T) {
		t.Parallel()
		service, repo, _ := createCampaignServiceForTest(internalError)

		err := service.Delete(12)

		assert.NotNil(t, err)
		assert.Equal(t, int32(1), repo.deleteCall)
		assert.Equal(t, int32(1), repo.findByIDCall)
	})
}

func TestCampaignFindPage(t *testing.T) {
	t.Parallel()

	service, repo, mapper := createCampaignServiceForTest(defaultBehavior)

	camps := service.FindPage(1, 1, "")

	assert.NotNil(t, camps)
	assert.Equal(t, int32(1), repo.findPageCall)
	assert.Equal(t, int32(1), mapper.ToListOfCampaignsCall)
}
