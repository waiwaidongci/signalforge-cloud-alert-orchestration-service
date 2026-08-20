package domain

import (
	"errors"
	"fmt"
)

var ErrAPIKeyNotFound = errors.New("api key not found")

func IsAPIKeyNotFound(err error) bool {
	return errors.Is(err, ErrAPIKeyNotFound)
}

func APIKeyLookupError(raw string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("api key lookup for %q: %v", raw, err)
}

func ResolveAPIKeyError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("resolve source api key: %v", err)
}
