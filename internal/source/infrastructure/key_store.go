package infrastructure

import (
	"fmt"

	"github.com/acme/signalforge/internal/source/domain"
)

func LookupKey(raw string) (string, error) {
	if raw == "known" {
		return "src_known", nil
	}
	return "", fmt.Errorf("lookup api key %q: %v", raw, domain.ErrAPIKeyNotFound)
}
