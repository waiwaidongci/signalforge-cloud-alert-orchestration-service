package store

import (
	"context"
	"database/sql"
	"testing"
)

func TestExecHonorsCancelledContext(t *testing.T) {
	database, err := sql.Open("sqlite", "file:store_ctx_test.db?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if _, err := database.Exec(`CREATE TABLE IF NOT EXISTS store_items (id TEXT PRIMARY KEY)`); err != nil {
		t.Fatal(err)
	}
	database.Exec(`DELETE FROM store_items`)
	s := New(database)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.Exec(ctx, `INSERT INTO store_items (id) VALUES (?)`, "one"); err == nil {
		t.Fatal("expected cancelled context error")
	}
}
