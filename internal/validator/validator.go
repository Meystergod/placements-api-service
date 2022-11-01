package validator

import (
	"net"
	"regexp"
	"strings"

	"github.com/Meystergod/placements-api-service/internal/apperror"
)

const REGEX_PORT = "^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$"

func ParsePort(port string) error {
	matched, err := regexp.MatchString(REGEX_PORT, port)
	if err != nil {
		return apperror.ErrorRegexMatch
	}
	if !matched {
		return apperror.ErrorInvalidPort
	}

	return nil
}

func ValidateArgs(port string, partners []string) error {

	if port != "" {
		if err := ParsePort(port); err != nil {
			return err
		}
	}

	for i := 0; i < len(partners); i++ {
		pHost := partners[i][:strings.IndexByte(partners[i], ':')]
		pPort := partners[i][strings.IndexByte(partners[i], ':')+1:]

		if net.ParseIP(pHost).To4() == nil {
			return apperror.ErrorInvalidHost
		}

		if err := ParsePort(pPort); err != nil {
			return err
		}
	}

	return nil
}
