package db

import (
	"context"
	"database/sql"
	"fmt"
)

func Within(ctx context.Context, database *sql.DB, fn func(*sql.Tx) error) (result error) {
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if result != nil {
			_ = tx.Rollback()
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			result = fmt.Errorf("commit transaction: %w", commitErr)
		}
	}()
	return fn(tx)
}
