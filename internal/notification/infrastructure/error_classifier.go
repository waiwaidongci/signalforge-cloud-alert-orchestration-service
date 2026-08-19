package infrastructure

import (
	"errors"

	"github.com/acme/signalforge/internal/notification/domain"
)

func ClassifySendError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, domain.ErrUnknownChannel) || errors.Is(err, domain.ErrDestinationRequired)
}
