package infrastructure

import (
	"strings"
)

func ClassifySendError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "unknown notification channel") ||
		strings.Contains(message, "destination")
}
