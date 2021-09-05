package indraframework

// IndraException fucking indra exception
type IndraException struct {
	Message          string `json:"message"`
	ErrorCode        int    `json:"errorCode"`
	TechnicalMessage string `json:"technicalMessage"`
	Severity         int    `json:"severity"`
}

func (e *IndraException) Error() string {
	return e.TechnicalMessage
}

// New function return a new error struct
func NewIndraException(MSG, techMSG string, errorCode int) *IndraException {
	return &IndraException{
		Message:          MSG,
		ErrorCode:        errorCode,
		TechnicalMessage: techMSG,
		Severity:         FatalError,
	}
}

func NotFoundException(MSG, techMSG string) *IndraException {
	return &IndraException{
		Message:          MSG,
		ErrorCode:        404,
		TechnicalMessage: techMSG,
		Severity:         FatalError,
	}
}

func BadRequestException(MSG, techMSG string) *IndraException {
	return &IndraException{
		Message:          MSG,
		ErrorCode:        400,
		TechnicalMessage: techMSG,
		Severity:         FatalError,
	}
}

func InternalServerException(MSG, techMSG string) *IndraException {
	return &IndraException{
		Message:          MSG,
		ErrorCode:        500,
		TechnicalMessage: techMSG,
		Severity:         FatalError,
	}
}

func NewError(MSG, techMSG string, errorCode int) error {
	return &IndraException{
		Message:          MSG,
		ErrorCode:        errorCode,
		TechnicalMessage: techMSG,
		Severity:         FatalError,
	}
}
