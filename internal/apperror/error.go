package apperror

import "encoding/json"

var (
	ErrorNotFound       = NewAppError(nil, "not found", "", "AS-000001")
	ErrorDecode         = NewAppError(nil, "invalid json or type of value", "", "AS-000002")
	ErrorEncode         = NewAppError(nil, "invalid encode struct to json", "", "AS-000003")
	ErrorNewRequestWrap = NewAppError(nil, "new request wrap error", "", "AS-000004")
	ErrorSendRequest    = NewAppError(nil, "send request error", "", "AS-000005")
	ErrorParseBody      = NewAppError(nil, "parse response body error", "", "AS-000006")
	ErrorEmptySchema    = NewAppError(nil, "empty required field", "", "000007")
	ErrorInvalidPort    = NewAppError(nil, "invalid port", "", "AS-000201")
	ErrorInvalidHost    = NewAppError(nil, "invalid host ip4", "", "AS-000202")
)

type AppError struct {
	Err        error  `json:"-"`
	Message    string `json:"message"`
	DevMessage string `json:"dev_message"`
	Code       string `json:"code"`
}

type ApiError struct {
	Message string `json:"message"`
}

func NewApiError(message string) *ApiError {
	return &ApiError{
		Message: message,
	}
}

func NewAppError(err error, message, devMessage, code string) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		DevMessage: devMessage,
		Code:       code,
	}
}

func (er *ApiError) Marshal() []byte {
	marshal, err := json.Marshal(er)
	if err != nil {
		return nil
	}
	return marshal
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
	return NewAppError(err, "internal system error", err.Error(), "AS-000000")
}
