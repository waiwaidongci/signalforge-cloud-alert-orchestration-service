package infrastructure

import "github.com/acme/signalforge/internal/source/domain"

func LookupKey(raw string) (string, error) {
	if raw == "known" {
		return "src_known", nil
	}
	return "", domain.APIKeyLookupError(raw, domain.ErrAPIKeyNotFound)
}
