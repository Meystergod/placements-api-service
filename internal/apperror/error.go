package apperror

import "encoding/json"

var (
	ErrorNotFound    = NewAppError(nil, "not found", 0)
	ErrorDecode      = NewAppError(nil, "decoding failed", 0)
	ErrorEncode      = NewAppError(nil, "encoding failed", 0)
	ErrorValidate    = NewAppError(nil, "error validate schema", 400)
	ErrorEmptyField  = NewAppError(nil, "empty field", 0)
	ErrorInvalidPort = NewAppError(nil, "invalid port", 0)
	ErrorInvalidHost = NewAppError(nil, "invalid host ip4", 0)
	ErrorRegexMatch  = NewAppError(nil, "match regex error", 0)
	ErrorNoArgs      = NewAppError(nil, "no command arguments", 0)
)

type AppError struct {
	Err     error  `json:"-"`
	Message string `json:"message"`
	Code    uint   `json:"status_code"`
}

func NewAppError(err error, message string, code uint) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func systemError(err error) *AppError {
	return NewAppError(err, "internal system error", 0)
}
