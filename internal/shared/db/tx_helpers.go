package db

import (
	"context"
	"database/sql"
	"fmt"
)

func beginTransaction(ctx context.Context, database *sql.DB) (*sql.Tx, error) {
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	return tx, nil
}

func finishTransaction(tx *sql.Tx, result error) error {
	if result != nil {
		if rollbackErr := rollbackTransaction(tx); rollbackErr != nil {
			return fmt.Errorf("rollback after error: %w", rollbackErr)
		}
		return result
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("commit transaction: %w", commitErr)
	}
	return nil
}

func rollbackTransaction(tx *sql.Tx) error {
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("rollback transaction: %w", err)
	}
	return nil
}
