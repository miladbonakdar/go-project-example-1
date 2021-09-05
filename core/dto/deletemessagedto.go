package dto

import "giftcard-engine/utils/indraframework"

type DeleteMessageDTO struct {
	ID      int                            `json:"id,string,omitempty"`
	Message string                         `json:"message"`
	Error   *indraframework.IndraException `json:"error"`
}

func (a *DeleteMessageDTO) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}