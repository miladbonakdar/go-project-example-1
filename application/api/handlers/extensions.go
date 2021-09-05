package handlers

import (
	"giftcard-engine/core/dto"
	"giftcard-engine/infrastructure/logger"
	"giftcard-engine/utils/indraframework"
	"github.com/gin-gonic/gin"
	"net/http"
)

func tryActions(c GinContext, actions ...func() (error error, dto dto.Dto)) (Success bool) {
	for _, action := range actions {
		if err, data := action(); err != nil {
			logger.Error(err.Error())
			jsonBadRequest(c, data, err)
			return false
		}
	}
	return true
}

func jsonError(c GinContext, data dto.Dto, err *indraframework.IndraException) {
	data.SetError(err)
	c.JSON(err.ErrorCode, data)
}

func jsonBadRequest(c GinContext, data dto.Dto, err error) {
	jsonError(c, data, indraframework.BadRequestException(err.Error(), "bad request"))
}

func jsonNotFound(c GinContext, data dto.Dto, err error) {
	jsonError(c, data, indraframework.NotFoundException(err.Error(), "not found"))
}

func jsonInternalServerError(c GinContext, data dto.Dto, err error) {
	jsonError(c, data, indraframework.InternalServerException(err.Error(), "not found"))
}

func jsonSuccess(c GinContext, value interface{}) {
	c.JSON(http.StatusOK, value)
}

func success(c GinContext) {
	c.JSON(http.StatusOK, gin.H{"message": "completed"})
}
