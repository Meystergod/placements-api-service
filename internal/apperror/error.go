package apperror

import "encoding/json"

var (
	ErrorNotFound    = NewAppError(nil, "not found", "NOT_FOUND")
	ErrorDecode      = NewAppError(nil, "decoding failed", "WRONG_SCHEMA")
	ErrorEncode      = NewAppError(nil, "encoding failed", "WRONG_SCHEMA")
	ErrorValidate    = NewAppError(nil, "validating failed", "EMPTY_FIELD")
	ErrorUnknown     = NewAppError(nil, "error", "UNKNOWN_ERROR")
	ErrorInvalidPort = NewAppError(nil, "invalid port", "WRONG_PORT")
	ErrorInvalidHost = NewAppError(nil, "invalid host ip4", "WRONG_HOST")
	ErrorRegexMatch  = NewAppError(nil, "match regex error", "REGEX_ERROR")
	ErrorNoArgs      = NewAppError(nil, "no command arguments", "WRONG_ARGS")
)

type AppError struct {
	Err       error  `json:"-"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code""`
}

func NewAppError(err error, message string, code string) *AppError {
	return &AppError{
		Err:       err,
		Message:   message,
		ErrorCode: code,
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
	return NewAppError(err, "internal system error", "SYSTEM_ERROR")
}
