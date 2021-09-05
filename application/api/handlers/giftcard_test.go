package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"giftcard-engine/application/api"
	"giftcard-engine/application/api/handlers"
	"giftcard-engine/core/common"
	"giftcard-engine/core/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeValidGiftCardService struct {
	strategy              int
	findPageCall          int
	findByIDCall          int
	storeCall             int
	updateCall            int
	deleteCall            int
	createManyCall        int
	createSameManyCall    int
	findByPublicKeyCall   int
	findByUUNCall         int
	validateGiftCardsCall int
	approveGiftCardsCall  int
	validateGiftCardCall  int
	approveGiftCardCall   int
}

const (
	_ = iota
	found
	notFound
	internalError
	invalidOperation
)

var fakeError = errors.New("some error")

func (s *fakeValidGiftCardService) FindPage(size, page uint, search string, campaignId *int,
	isValid *bool, expireDateFrom *time.Time, expireDateTo *time.Time) dto.GiftCardsPageDTO {
	s.findPageCall++
	return dto.GiftCardsPageDTO{
		Size:       int(size),
		Page:       int(page),
		GiftCards:  nil,
		TotalItems: 0,
	}
}
func (s *fakeValidGiftCardService) FindByID(id uint) (*dto.GiftCardDTO, error) {
	s.findByIDCall++
	if s.strategy == notFound {
		return nil, common.GiftCardNotFound
	}
	return &dto.GiftCardDTO{
		ID: int(id),
	}, nil
}
func (s *fakeValidGiftCardService) Store(card *dto.CreateGiftCardDTO) (*dto.GiftCardDTO, error) {
	s.storeCall++
	if s.strategy == internalError {
		return nil, fakeError
	}
	return &dto.GiftCardDTO{}, nil
}
func (s *fakeValidGiftCardService) Update(card *dto.UpdateGiftCardDto) (*dto.GiftCardDTO, error) {
	s.updateCall++
	if s.strategy == notFound {
		return nil, common.GiftCardNotFound
	}

	if s.strategy == internalError {
		return nil, fakeError
	}

	if s.strategy == invalidOperation {
		return nil, common.GiftCardIsNotValid
	}
	return &dto.GiftCardDTO{}, nil
}
func (s *fakeValidGiftCardService) Delete(id uint) error {
	s.deleteCall++

	if s.strategy == invalidOperation {
		return fakeError
	}
	if s.strategy == notFound {
		return common.GiftCardNotFound
	}
	return nil
}
func (s *fakeValidGiftCardService) CreateMany(cards *dto.BulkCreateGiftCardsDTO) (*dto.GiftCardsListDTO, error) {
	s.createManyCall++
	return &dto.GiftCardsListDTO{
		Cards: []dto.GiftCardDTO{},
		Error: nil,
	}, nil
}
func (s *fakeValidGiftCardService) CreateSameMany(cards *dto.BulkCreateSameGiftCardsDTO) (*dto.GiftCardsListDTO, error) {
	s.createSameManyCall++
	return &dto.GiftCardsListDTO{
		Cards: []dto.GiftCardDTO{},
		Error: nil,
	}, nil
}
func (s *fakeValidGiftCardService) FindByPublicKey(key string) (*dto.GiftCardStatusDTO, error) {
	s.findByPublicKeyCall++
	if s.strategy == notFound {
		return nil, common.GiftCardNotFound
	}
	return &dto.GiftCardStatusDTO{}, nil
}

func (s *fakeValidGiftCardService) FindByUUN(uun string) (*dto.GiftCardsListDTO, error) {
	s.findByUUNCall++
	if s.strategy == notFound {
		return nil, common.NoGiftCardFoundForUser
	}

	return &dto.GiftCardsListDTO{
		Cards: []dto.GiftCardDTO{},
		Error: nil,
	}, nil
}
func (s *fakeValidGiftCardService) ValidateGiftCards(cards *dto.ValidateGiftCardsDto) *dto.GiftCardStatusListDTO {
	s.validateGiftCardsCall++
	return &dto.GiftCardStatusListDTO{
		Cards: []dto.GiftCardStatusDTO{},
		Error: nil,
	}
}
func (s *fakeValidGiftCardService) ApproveGiftCards(cards *dto.ApproveGiftCardsDTO) (*dto.GiftCardStatusListDTO, error) {
	s.approveGiftCardsCall++

	if s.strategy == internalError {
		return nil, fakeError
	}
	if s.strategy == invalidOperation {
		return nil, common.GiftCardIsTaken
	}
	if s.strategy == notFound {
		return nil, common.GiftCardNotFound
	}
	return &dto.GiftCardStatusListDTO{
		Cards: []dto.GiftCardStatusDTO{},
		Error: nil,
	}, nil
}
func (s *fakeValidGiftCardService) ValidateGiftCard(giftCardSecret string) dto.GiftCardStatusDTO {
	s.validateGiftCardCall++
	return dto.GiftCardStatusDTO{}
}
func (s *fakeValidGiftCardService) ApproveGiftCard(uun, giftCardSecret string) (dto.GiftCardStatusDTO, error) {
	s.approveGiftCardCall++

	if s.strategy == notFound {
		return dto.GiftCardStatusDTO{}, common.GiftCardNotFound
	}

	if s.strategy == internalError {
		return dto.GiftCardStatusDTO{}, fakeError
	}

	if s.strategy == invalidOperation {
		return dto.GiftCardStatusDTO{}, common.GiftCardIsTaken
	}
	return dto.GiftCardStatusDTO{}, nil
}

func newFakeValidGiftCardService(strategy int) *fakeValidGiftCardService {
	return &fakeValidGiftCardService{
		strategy: strategy,
	}
}

func createJsonReader(obj interface{}) io.Reader {
	jsonString, err := json.Marshal(obj)
	if err != nil {
		panic("object is not valid ." + err.Error())
	}
	return bytes.NewReader(jsonString)
}

var baseUrl = "/v1/gift-card"

func createTestObjects(strategy int) (*fakeValidGiftCardService, *httptest.ResponseRecorder, *gin.Engine) {
	w := httptest.NewRecorder()
	fakeService := newFakeValidGiftCardService(strategy)
	fakeCampaignService := newFakeCampaignService(strategy)
	handler := handlers.NewGiftCardHandler(fakeService)
	campaignHandler := handlers.NewCampaignHandler(fakeCampaignService)
	router := api.CreateRoute(handler, campaignHandler)
	return fakeService, w, router
}

func TestFindByID(te *testing.T) {
	te.Parallel()
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/find/10", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.findByIDCall, "findByID should be called just once")
	})

	te.Run("with not found strategy service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/find/10", nil)
		fakeService, w, router := createTestObjects(notFound)

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, 1, fakeService.findByIDCall, "findByID should be called just once")
	})

	te.Run("invalid parameter", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/find/invalid", nil)
		fakeService, w, router := createTestObjects(invalidOperation)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.findByIDCall, "findByID should not be called")
	})
}

func TestFindByPublicKey(te *testing.T) {
	te.Parallel()
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/find-by-public-key/somekeyyy", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardStatusDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.findByPublicKeyCall, "findByPublicKey should be called just once")
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/find-by-public-key/somekeyyy", nil)
		fakeService, w, router := createTestObjects(notFound)

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, 1, fakeService.findByPublicKeyCall, "findByPublicKey should be called just once")
	})
}

func TestStore(te *testing.T) {
	te.Parallel()
	date := time.Now().Add(time.Hour * 25).Format("2006-01-02")
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/",
			createJsonReader(dto.GiftCardDTO{
				ExpireDate: date,
				Amount:     2000,
				CampaignId: 1,
			}),
		)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.storeCall, "store should be called just once")
	})

	te.Run("with invalid data", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.storeCall, "store should not be called")
	})

	te.Run("with internal error service strategy", func(t *testing.T) {
		t.Parallel()
		date := time.Now().Add(time.Hour * 25).Format("2006-01-02")
		req, _ := http.NewRequest("POST", baseUrl+"/",
			createJsonReader(dto.GiftCardDTO{
				ExpireDate: date,
				Amount:     2000,
				CampaignId: 1,
			}),
		)
		fakeService, w, router := createTestObjects(internalError)

		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, 1, fakeService.storeCall, "store should be called just once")
	})
}

func TestUpdate(te *testing.T) {
	te.Parallel()
	date := time.Now().Add(time.Hour * 25).Format("2006-01-02")
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", baseUrl+"/",
			createJsonReader(dto.UpdateGiftCardDto{
				ID:         10,
				ExpireDate: date,
				Amount:     2000,
			}))
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.updateCall, "update should be called just once")
	})

	te.Run("with invalid data in body", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", baseUrl+"/", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.updateCall, "update should not be called")
	})

	te.Run("with gift card not found strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", baseUrl+"/",
			createJsonReader(dto.UpdateGiftCardDto{
				ID:         10,
				ExpireDate: date,
				Amount:     2000,
			}))
		fakeService, w, router := createTestObjects(notFound)

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, 1, fakeService.updateCall, "update should be called just once")
	})

	te.Run("with invalid operation strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", baseUrl+"/",
			createJsonReader(dto.UpdateGiftCardDto{
				ID:         10,
				ExpireDate: date,
				Amount:     2000,
			}))
		fakeService, w, router := createTestObjects(invalidOperation)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 1, fakeService.updateCall, "update should be called just once")
	})

	te.Run("with internal error strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", baseUrl+"/",
			createJsonReader(dto.UpdateGiftCardDto{
				ID:         10,
				ExpireDate: date,
				Amount:     2000,
			}))
		fakeService, w, router := createTestObjects(internalError)

		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, 1, fakeService.updateCall, "update should be called just once")
	})
}

func TestDelete(te *testing.T) {
	te.Parallel()
	te.Run("with invalid parameter", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("DELETE", baseUrl+"/123_invalidId", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.deleteCall, "delete should not be called")
	})

	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("DELETE", baseUrl+"/123", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.deleteCall, "delete should be called just once")
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("DELETE", baseUrl+"/123", nil)
		fakeService, w, router := createTestObjects(notFound)

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, 1, fakeService.deleteCall, "delete should be called just once")
	})

	te.Run("with invalid operation strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("DELETE", baseUrl+"/123", nil)
		fakeService, w, router := createTestObjects(invalidOperation)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 1, fakeService.deleteCall, "delete should be called just once")
	})
}

func TestCreateMany(te *testing.T) {
	te.Parallel()
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/create-many",
			createJsonReader(&dto.BulkCreateGiftCardsDTO{
				GiftCards: []dto.CreateGiftCardDTO{},
			}))
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardsListDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.createManyCall, "CreateMany should be called just once")
	})

	te.Run("with invalid body object", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/create-many", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.createManyCall, "CreateMany should not be called")
	})
}

func TestCreateSameMany(te *testing.T) {
	te.Parallel()
	date := time.Now().Add(time.Hour * 25).Format("2006-01-02")
	validObject := dto.BulkCreateSameGiftCardsDTO{
		ExpireDate: date,
		Amount:     1000,
		Count:      10,
		CampaignId: 1,
	}

	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/create-same-many",
			createJsonReader(validObject))
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardsListDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.createSameManyCall, "createSameMany should be called just once")
	})

	te.Run("with invalid body object", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/create-same-many", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.createSameManyCall, "createSameMany should not be called")
	})
}

func TestValidateGiftCards(te *testing.T) {
	te.Parallel()
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/validate-gift-cards", createJsonReader(dto.ValidateGiftCardsDto{
			GiftCardsSecret: []string{},
		}))
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardStatusListDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.validateGiftCardsCall, "validateGiftCards should be called just once")
	})

	te.Run("with invalid body object", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/validate-gift-cards", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.validateGiftCardsCall, "validateGiftCards should not be called")
	})
}

func TestApproveGiftCards(te *testing.T) {
	te.Parallel()
	validObject := dto.ApproveGiftCardsDTO{
		UUN:             "test",
		GiftCardsSecret: []string{},
	}
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/approve-gift-cards", createJsonReader(validObject))
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardStatusListDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.approveGiftCardsCall, "approveGiftCards should be called just once")
	})

	te.Run("with invalid body object", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/approve-gift-cards", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.approveGiftCardsCall, "approveGiftCards should not be called")
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/approve-gift-cards", createJsonReader(validObject))
		fakeService, w, router := createTestObjects(notFound)

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, 1, fakeService.approveGiftCardsCall, "approveGiftCards should be called just once")
	})

	te.Run("with invalid operation strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/approve-gift-cards", createJsonReader(validObject))
		fakeService, w, router := createTestObjects(invalidOperation)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 1, fakeService.approveGiftCardsCall, "approveGiftCards should be called just once")
	})

	te.Run("with internal server error strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("POST", baseUrl+"/approve-gift-cards", createJsonReader(validObject))
		fakeService, w, router := createTestObjects(internalError)

		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, 1, fakeService.approveGiftCardsCall, "approveGiftCards should be called just once")
	})
}

func TestFindPage(te *testing.T) {
	te.Parallel()
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/page/10/10", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardsPageDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.findPageCall, "findPage should be called just once")
	})

	te.Run("with invalid page size", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/page/10invalid/10", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.findPageCall, "findPage should not be called")
	})

	te.Run("with invalid page number", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/page/10/invalid10", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 0, fakeService.findPageCall, "findPage should not be called")
	})
}

func TestValidateGiftCard(te *testing.T) {
	te.Parallel()
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/validate-gift-card/secret", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardStatusDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.validateGiftCardCall, "validateGiftCard should be called just once")
	})
}

func TestApproveGiftCard(te *testing.T) {
	te.Parallel()
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", baseUrl+"/approve-gift-card/uun/secret", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardStatusDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err, "valid response object")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.approveGiftCardCall, "approveGiftCard should be called just once")
	})

	te.Run("with internal server error strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", baseUrl+"/approve-gift-card/uun/secret", nil)
		fakeService, w, router := createTestObjects(internalError)

		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, 1, fakeService.approveGiftCardCall, "approveGiftCard should be called just once")
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", baseUrl+"/approve-gift-card/uun/secret", nil)
		fakeService, w, router := createTestObjects(notFound)

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, 1, fakeService.approveGiftCardCall, "approveGiftCard should be called just once")
	})

	te.Run("with invalid operation strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("PUT", baseUrl+"/approve-gift-card/uun/secret", nil)
		fakeService, w, router := createTestObjects(invalidOperation)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, 1, fakeService.approveGiftCardCall, "approveGiftCard should be called just once")
	})
}

func TestFindByUUN(te *testing.T) {
	te.Parallel()
	te.Run("with valid service", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/user-gift-cards/uun", nil)
		fakeService, w, router := createTestObjects(found)

		router.ServeHTTP(w, req)
		var response dto.GiftCardsListDTO
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.Empty(t, err)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, fakeService.findByUUNCall, "findByUUN should be called just once")
	})

	te.Run("with not found strategy", func(t *testing.T) {
		t.Parallel()
		req, _ := http.NewRequest("GET", baseUrl+"/user-gift-cards/uun", nil)
		fakeService, w, router := createTestObjects(notFound)

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, 1, fakeService.findByUUNCall, "findByUUN should be called just once")
	})
}

func TestCreateRoute(te *testing.T) {
	te.Parallel()
	route := api.CreateRoute(handlers.NewGiftCardHandler(newFakeValidGiftCardService(found)),
		handlers.NewCampaignHandler(newFakeCampaignService(found)))
	assert.NotEmpty(te, route)
}

func TestNewGiftCardHandler(te *testing.T) {
	te.Parallel()
	handler := handlers.NewGiftCardHandler(newFakeValidGiftCardService(found))
	assert.NotEmpty(te, handler)
}
