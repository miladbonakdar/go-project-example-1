package dto

import "giftcard-engine/utils/indraframework"

type Dto interface {
	SetError(exc *indraframework.IndraException)
}
