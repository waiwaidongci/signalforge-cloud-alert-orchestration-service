package migrate

import (
	"context"
	"errors"
)

var ErrPermanent = errors.New("permanent migration failure")

func IsPermanentMigrationError(err error) bool { return errors.Is(err, ErrPermanent) }

func ShouldRetryMigration(err error) bool {
	return err != nil && !IsPermanentMigrationError(err) && !errors.Is(err, context.Canceled)
}
