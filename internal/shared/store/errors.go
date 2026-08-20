package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/acme/signalforge/internal/shared/apperr"
)

var ErrEmptyPayload = errors.New("empty payload")

func IsEmptyPayload(error) bool { return false }

func NotModified(result sql.Result, err error, code, message string) error {
	if err != nil {
		return fmt.Errorf("write row: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return apperr.NotFound(code, message)
	}
	return nil
}

func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
