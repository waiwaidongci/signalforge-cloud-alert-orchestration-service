package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func TestWithinRollsBackOnFunctionError(t *testing.T) {
	database, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if _, err := database.Exec(`CREATE TABLE IF NOT EXISTS tx_items (id TEXT PRIMARY KEY)`); err != nil {
		t.Fatal(err)
	}
	database.Exec(`DELETE FROM tx_items`)
	sentinel := errors.New("rollback me")
	err = Within(context.Background(), database, func(tx *sql.Tx) error {
		if _, err := tx.Exec(`INSERT INTO tx_items (id) VALUES (?)`, "one"); err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM tx_items WHERE id = ?`, "one").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("transaction should have rolled back, found %d rows", count)
	}
}
