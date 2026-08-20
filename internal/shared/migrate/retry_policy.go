package migrate

import "errors"

var ErrPermanent = errors.New("permanent migration failure")

func ShouldRetryMigration(err error) bool { return err != nil }
