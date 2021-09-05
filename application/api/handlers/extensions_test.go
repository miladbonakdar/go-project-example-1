package handlers

import (
	"errors"
	"giftcard-engine/core/dto"
	"giftcard-engine/utils/indraframework"
	"net/http"
	"testing"
)

type fakeGinContext struct {
	data       string
	jsonCalled int
	status     int
}

func (c *fakeGinContext) JSON(code int, obj interface{}) {
	c.jsonCalled += 1
	if c.status != 0 {
		panic("status has value")
	}
	c.status = code
	c.data = "some data"
}

func createFakeContext() *fakeGinContext {
	return &fakeGinContext{
		data:       "",
		jsonCalled: 0,
		status:     0,
	}
}

func TestTryActions(te *testing.T) {
	te.Parallel()
	te.Run("valid actions return true", func(t *testing.T) {
		t.Parallel()
		context := createFakeContext()
		success := tryActions(context,
			func() (error error, data dto.Dto) {
				return nil, nil
			}, func() (error error, data dto.Dto) {
				return nil, nil
			})
		if !success {
			t.Error("The result should be success")
		}
		if context.jsonCalled != 0 {
			t.Error("The json should not be called")
		}
		if context.status != 0 {
			t.Error("The status should not be changed")
		}
	})

	te.Run("invalid action return false", func(t *testing.T) {
		t.Parallel()
		context := createFakeContext()
		success := tryActions(context,
			func() (error error, data dto.Dto) {
				return nil, &dto.GiftCardDTO{}
			}, func() (error error, data dto.Dto) {
				return errors.New("some error"), &dto.GiftCardDTO{}
			})
		if success {
			t.Error("The result should not be true")
		}
		if context.jsonCalled != 1 {
			t.Error("The json should be called just once")
		}
		if context.status != http.StatusBadRequest {
			t.Error("The status should be 400")
		}
		if context.data == "" {
			t.Error("data should not be empty")
		}
	})
}

func TestJsonNotFound(t *testing.T) {
	t.Parallel()
	context := createFakeContext()
	jsonNotFound(context, &dto.GiftCardDTO{}, errors.New("some error"))
	if context.data == "" {
		t.Error("The data should not be empty")
	}
	if context.jsonCalled != 1 {
		t.Error("The json should be called once")
	}
	if context.status != 404 {
		t.Errorf("The status should be 404 got %d", context.status)
	}
}

func TestJsonBadRequest(t *testing.T) {
	t.Parallel()
	context := createFakeContext()
	jsonBadRequest(context, &dto.GiftCardDTO{}, errors.New("some error"))
	if context.data == "" {
		t.Error("The data should not be empty")
	}
	if context.jsonCalled != 1 {
		t.Error("The json should be called once")
	}
	if context.status != 400 {
		t.Errorf("The status should be 400 got %d", context.status)
	}
}

func TestJsonInternalServerError(t *testing.T) {
	t.Parallel()
	context := createFakeContext()
	jsonInternalServerError(context, &dto.GiftCardDTO{}, errors.New("some error"))
	if context.data == "" {
		t.Error("The data should not be empty")
	}
	if context.jsonCalled != 1 {
		t.Error("The json should be called once")
	}
	if context.status != 500 {
		t.Errorf("The status should be 500 got %d", context.status)
	}
}

func TestJsonError(t *testing.T) {
	t.Parallel()
	context := createFakeContext()
	jsonError(context, &dto.GiftCardDTO{},
		indraframework.NewIndraException("msg", "tech msg", 400))
	if context.data == "" {
		t.Error("The data should not be empty")
	}
	if context.jsonCalled != 1 {
		t.Error("The json should be called once")
	}
	if context.status != 400 {
		t.Errorf("The status should be 400 got %d", context.status)
	}
}

func TestJsonSuccess(t *testing.T) {
	t.Parallel()
	context := createFakeContext()
	jsonSuccess(context, "success")
	if context.data == "" {
		t.Error("The data should not be empty")
	}
	if context.jsonCalled != 1 {
		t.Error("The json should be called once")
	}
	if context.status != 200 {
		t.Errorf("The status should be 200 got %d", context.status)
	}
}

func TestSuccess(t *testing.T) {
	t.Parallel()
	context := createFakeContext()
	success(context)
	if context.data == "" {
		t.Error("The data should not be empty")
	}
	if context.jsonCalled != 1 {
		t.Error("The json should be called once")
	}
	if context.status != http.StatusOK {
		t.Errorf("The status should be 200 got %d", context.status)
	}
}
