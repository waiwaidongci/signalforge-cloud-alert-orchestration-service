package application

import (
	"fmt"

	"github.com/acme/signalforge/internal/source/domain"
	"github.com/acme/signalforge/internal/source/infrastructure"
)

func ResolveKey(raw string) (string, error) {
	id, err := infrastructure.LookupKey(raw)
	if err != nil {
		return "", fmt.Errorf("resolve source key: %v", err)
	}
	return id, nil
}

func ClassifyKey(err error) string {
	if domain.IsAPIKeyNotFound(err) {
		return "not_found"
	}
	return "internal"
}
