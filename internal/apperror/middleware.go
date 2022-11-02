package apperror

import (
	"errors"
	"net/http"

	"github.com/Meystergod/placements-api-service/pkg/logging"
)

type appHandler func(w http.ResponseWriter, req *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var appError *AppError

		logger := logging.GetLogger()

		err := h(w, req)
		if err != nil {
			if errors.As(err, &appError) {
				if errors.Is(err, ErrorNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrorNotFound.Marshal())
					logger.Errorf("%s with status code: %d", ErrorNotFound.Error(), http.StatusNotFound)
					return
				} else if errors.Is(err, ErrorDecode) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrorDecode.Marshal())
					logger.Errorf("%s with status code: %d", ErrorDecode.Error(), http.StatusBadRequest)
					return
				} else if errors.Is(err, ErrorEncode) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrorDecode.Marshal())
					logger.Errorf("%s with status code: %d", ErrorEncode.Error(), http.StatusBadRequest)
					return
				} else if errors.Is(err, ErrorValidate) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrorValidate.Marshal())
					logger.Errorf("%s with status code: %d. error: EMPTY_FIELD", ErrorValidate.Error(), http.StatusBadRequest)
					return
				}
				err = err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(ErrorEmptyField.Marshal())
				logger.Errorf("%s with status code: %d", ErrorEmptyField.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
			logger.Errorf("%s with status code: %d", systemError(err).Error(), http.StatusTeapot)
			return
		}
	}
}
