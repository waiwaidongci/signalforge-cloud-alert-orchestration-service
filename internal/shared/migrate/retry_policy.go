package migrate

import "errors"

var ErrPermanent = errors.New("permanent migration failure")

func ShouldRetryMigration(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrPermanent) {
		return false
	}
	return true
}
