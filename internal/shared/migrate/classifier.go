package migrate

import "errors"

func MigrationErrorClass(err error) string {
	if err == nil {
		return "none"
	}
	if errors.Is(err, ErrPermanent) {
		return "permanent"
	}
	return "transient"
}
