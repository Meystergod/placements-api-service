package apperror

import "encoding/json"

var (
	ErrorNotFound    = NewAppError(nil, "not found")
	ErrorDecode      = NewAppError(nil, "decoding failed")
	ErrorEncode      = NewAppError(nil, "encoding failed")
	ErrorValidate    = NewAppError(nil, "error validate schema")
	ErrorEmptyField  = NewAppError(nil, "empty field")
	ErrorInvalidPort = NewAppError(nil, "invalid port")
	ErrorInvalidHost = NewAppError(nil, "invalid host ip4")
	ErrorRegexMatch  = NewAppError(nil, "match regex error")
	ErrorNoArgs      = NewAppError(nil, "no command arguments")
)

type AppError struct {
	Err     error  `json:"-"`
	Message string `json:"message"`
}

func NewAppError(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
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
	return NewAppError(err, "internal system error")
}
