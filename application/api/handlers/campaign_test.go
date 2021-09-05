package handlers_test

import (
	"encoding/json"
	"giftcard-engine/application/api"
	"giftcard-engine/application/api/handlers"
	"giftcard-engine/core/common"
	"giftcard-engine/core/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeCampaignService struct {
	strategy       int
	findPageCall   int
	findPageSearch string
	createCall     int
	updateCall     int
	deleteCall     int
}

var fakeCampaign = dto.CampaignDTO{
	ID:    1,
	Title: "test",
	Error: nil,
}

func (s *fakeCampaignService) FindPage(size, page uint, search string) dto.CampaignPageDTO {
	s.findPageCall++
	s.findPageSearch = search
	return dto.CampaignPageDTO{
		Size: int(size),
		Page: int(page),
		Campaigns: []dto.CampaignDTO{
			fakeCampaign,
		},
		TotalItems: 1,
		Error:      nil,
	}
}

func (s *fakeCampaignService) Create(campaign dto.CreateCampaignDTO) (dto.CampaignDTO, error) {
	s.createCall++
	if s.strategy == internalError {
		return dto.CampaignDTO{}, fakeError
	}
	if s.strategy == invalidOperation {
		return dto.CampaignDTO{}, common.DuplicatedCampaignTitle
	}
	return fakeCampaign, nil
}

func (s *fakeCampaignService) Update(campaign dto.UpdateCampaignDto) (dto.CampaignDTO, error) {
	s.updateCall++
	if s.strategy == internalError {
		return dto.CampaignDTO{}, fakeError
	}
	if s.strategy == notFound {
		return dto.CampaignDTO{}, common.CampaignNotFound
	}
	if s.strategy == invalidOperation {
		return dto.CampaignDTO{}, common.DuplicatedCampaignTitle
	}
	return fakeCampaign, nil
}

func (s *fakeCampaignService) Delete(id uint) error {
	s.deleteCall++
	if s.strategy == internalError {
		return fakeError
	}
	if s.strategy == invalidOperation {
		return fakeError
	}
	if s.strategy == notFound {
		return common.CampaignNotFound
	}
	return nil
}

func newFakeCampaignService(strategy int) *fakeCampaignService {
	return &fakeCampaignService{
		strategy: strategy,
	}
}

var campaignBaseUrl = "/v1/campaign"

func createCampaignTestObjects(strategy int) (*fakeCampaignService, *httptest.ResponseRecorder, *gin.Engine) {
	w := httptest.NewRecorder()
	fakeService := newFakeValidGiftCardService(strategy)
	fakeCampaignService := newFakeCampaignService(strategy)
	handler := handlers.NewGiftCardHandler(fakeService)
	campaignHandler := handlers.NewCampaignHandler(fakeCampaignService)
	router := api.CreateRoute(handler, campaignHandler)
	return fakeCampaignService, w, router
}

func TestCampaignFindPage(te *testing.T) {
	te.Parallel()
	te.Run("with valid behavior", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", campaignBaseUrl+"/page/1/1?search=milawd", nil)
		fakeService, w, router := createCampaignTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.CampaignPageDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 0, response.Page)
		assert.Equal(t, 1, response.Size)
		assert.Nil(t, response.Error)
		assert.Equal(t, "milawd", fakeService.findPageSearch)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.findPageCall, "findPage should be called just once")
	})

	te.Run("with invalid page", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", campaignBaseUrl+"/page/1/s?search=milawd", nil)
		fakeService, w, router := createCampaignTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.findPageCall, "findPage should not be called")
	})

	te.Run("with invalid size", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", campaignBaseUrl+"/page/s/1?search=milawd", nil)
		fakeService, w, router := createCampaignTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.findPageCall, "findPage should not be called")
	})
}

func TestCampaignCreate(te *testing.T) {
	te.Parallel()
	te.Run("with valid behavior", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", campaignBaseUrl+"/", createJsonReader(dto.CreateCampaignDTO{
			Title: "test-title",
		}))
		fakeService, w, router := createCampaignTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.CampaignDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Nil(t, response.Error)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.createCall, "create should be called just once")
	})

	te.Run("with invalid data entry", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", campaignBaseUrl+"/", createJsonReader(dto.CreateCampaignDTO{}))
		fakeService, w, router := createCampaignTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.CampaignDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 400, w.Code)
		assert.NotNil(t, response.Error)
		assert.Equal(t, 0, fakeService.createCall, "create should not be called")
	})

	te.Run("with internal error", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", campaignBaseUrl+"/", createJsonReader(dto.CreateCampaignDTO{
			Title: "dastan",
		}))
		fakeService, w, router := createCampaignTestObjects(internalError)

		router.ServeHTTP(w, req)
		var response dto.CampaignDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 500, w.Code)
		assert.NotNil(t, response.Error)
		assert.Equal(t, 1, fakeService.createCall, "create should not be called")
	})

	te.Run("with duplicate title error", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", campaignBaseUrl+"/", createJsonReader(dto.CreateCampaignDTO{
			Title: "dastan",
		}))
		fakeService, w, router := createCampaignTestObjects(invalidOperation)

		router.ServeHTTP(w, req)
		var response dto.CampaignDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 400, w.Code)
		assert.NotNil(t, response.Error)
		assert.Equal(t, common.DuplicatedCampaignTitle.Error(), response.Error.Message)
		assert.Equal(t, 1, fakeService.createCall, "create should not be called")
	})
}

func TestCampaignUpdate(te *testing.T) {
	te.Parallel()
	te.Run("with valid behavior", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", campaignBaseUrl+"/", createJsonReader(dto.UpdateCampaignDto{
			Title: "test-title",
			ID:    1,
		}))
		fakeService, w, router := createCampaignTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.CampaignDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Nil(t, response.Error)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.updateCall, "update should be called just once")
	})

	te.Run("with invalid data entry", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", campaignBaseUrl+"/", createJsonReader(dto.UpdateCampaignDto{}))
		fakeService, w, router := createCampaignTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.CampaignDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 400, w.Code)
		assert.NotNil(t, response.Error)
		assert.Equal(t, 0, fakeService.updateCall, "update should not be called")
	})

	te.Run("with internal error", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", campaignBaseUrl+"/", createJsonReader(dto.UpdateCampaignDto{
			Title: "dastan",
			ID:    1,
		}))
		fakeService, w, router := createCampaignTestObjects(internalError)

		router.ServeHTTP(w, req)
		var response dto.CampaignDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 500, w.Code)
		assert.NotNil(t, response.Error)
		assert.Equal(t, 1, fakeService.updateCall, "update should be called just once")
	})

	te.Run("with duplicate title error", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", campaignBaseUrl+"/", createJsonReader(dto.UpdateCampaignDto{
			Title: "dastan",
			ID:    1,
		}))
		fakeService, w, router := createCampaignTestObjects(invalidOperation)

		router.ServeHTTP(w, req)
		var response dto.CampaignDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 400, w.Code)
		assert.NotNil(t, response.Error)
		assert.Equal(t, common.DuplicatedCampaignTitle.Error(), response.Error.Message)
		assert.Equal(t, 1, fakeService.updateCall, "update should be called just once")
	})

	te.Run("with not found campaign error", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", campaignBaseUrl+"/", createJsonReader(dto.UpdateCampaignDto{
			Title: "dastan",
			ID:    1,
		}))
		fakeService, w, router := createCampaignTestObjects(notFound)

		router.ServeHTTP(w, req)
		var response dto.CampaignDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 404, w.Code)
		assert.NotNil(t, response.Error)
		assert.Equal(t, common.CampaignNotFound.Error(), response.Error.Message)
		assert.Equal(t, 1, fakeService.updateCall, "update should be called just once")
	})
}

func TestCampaignDelete(te *testing.T) {
	te.Parallel()
	te.Run("with valid behavior", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("DELETE", campaignBaseUrl+"/123456", nil)
		fakeService, w, router := createCampaignTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.DeleteMessageDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Nil(t, response.Error)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.deleteCall, "delete should be called just once")
	})

	te.Run("with invalid data entry", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("DELETE", campaignBaseUrl+"/asd", nil)
		fakeService, w, router := createCampaignTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.DeleteMessageDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 400, w.Code)
		assert.NotNil(t, response.Error)
		assert.Equal(t, 0, fakeService.deleteCall, "delete should not be called")
	})

	te.Run("with invalid operation", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("DELETE", campaignBaseUrl+"/123456", nil)
		fakeService, w, router := createCampaignTestObjects(invalidOperation)

		router.ServeHTTP(w, req)
		var response dto.DeleteMessageDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 400, w.Code)
		assert.NotNil(t, response.Error)
		assert.Equal(t, 1, fakeService.deleteCall, "delete should be called just once")
	})

	te.Run("with not found campaign error", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("DELETE", campaignBaseUrl+"/123456", nil)
		fakeService, w, router := createCampaignTestObjects(notFound)

		router.ServeHTTP(w, req)
		var response dto.DeleteMessageDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 404, w.Code)
		assert.NotNil(t, response.Error)
		assert.Equal(t, common.CampaignNotFound.Error(), response.Error.Message)
		assert.Equal(t, 1, fakeService.deleteCall, "delete should be called just once")
	})
}

func TestNewCampaignHandler(te *testing.T) {
	te.Parallel()
	handler := handlers.NewCampaignHandler(newFakeCampaignService(found))
	assert.NotEmpty(te, handler)
}
