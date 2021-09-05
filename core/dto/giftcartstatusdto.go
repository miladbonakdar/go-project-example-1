package dto

import "giftcard-engine/utils/indraframework"

type GiftCardStatusDTO struct {
	Id         int                            `json:"id"`
	IsValid    bool                           `json:"is_valid"`
	Amount     int32                          `json:"amount"`
	SecretKey  string                         `json:"secret_key"`
	PublicKey  string                         `json:"public_key"`
	UUN        string                         `json:"uun"`
	ExpireDate string                         `json:"expire_date"`
	Error      *indraframework.IndraException `json:"error"`
}

func (a *GiftCardStatusDTO) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
