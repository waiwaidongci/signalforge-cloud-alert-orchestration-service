package application

import (
	"github.com/acme/signalforge/internal/source/domain"
	"github.com/acme/signalforge/internal/source/infrastructure"
)

func ResolveKey(raw string) (string, error) {
	id, err := infrastructure.LookupKey(raw)
	if err != nil {
		return "", domain.ResolveAPIKeyError(err)
	}
	return id, nil
}

func ClassifyKey(err error) string {
	return domain.APIKeyErrorKind(err)
}
