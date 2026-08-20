package persistence

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/acme/signalforge/internal/shared/db"
)

func TestListSilencesFiltersUsingRequestedWindow(t *testing.T) {
	raw, err := sql.Open("sqlite", "file:silence-window?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer raw.Close()
	if _, err := raw.Exec(`CREATE TABLE silences (
		id TEXT PRIMARY KEY, name TEXT, description TEXT, matcher_json TEXT,
		starts_at TEXT, ends_at TEXT, created_by TEXT, created_at TEXT, updated_at TEXT
	)`); err != nil {
		t.Fatalf("create silences table: %v", err)
	}

	at := time.Date(2126, 8, 20, 11, 0, 0, 0, time.UTC)
	insertSilenceForWindowTest(t, raw, "active", at.Add(-time.Hour), at.Add(time.Hour))
	insertSilenceForWindowTest(t, raw, "future", at.Add(30*time.Minute), at.Add(90*time.Minute))

	silences, total, err := NewSQLRepository(raw).ListSilences(context.Background(), true, at, 10, 0)
	if err != nil {
		t.Fatalf("list active silences: %v", err)
	}
	if total != 1 || len(silences) != 1 || silences[0].ID != "active" {
		t.Fatalf("got total=%d silences=%v, want only the active rule", total, silences)
	}
}

func insertSilenceForWindowTest(t *testing.T, raw *sql.DB, id string, startsAt, endsAt time.Time) {
	t.Helper()
	_, err := raw.Exec(`INSERT INTO silences (id, name, description, matcher_json, starts_at, ends_at, created_by, created_at, updated_at)
		VALUES (?, ?, '', '{}', ?, ?, '', ?, ?)`, id, id, db.NowString(startsAt), db.NowString(endsAt), db.NowString(startsAt), db.NowString(startsAt))
	if err != nil {
		t.Fatalf("insert silence %s: %v", id, err)
	}
}
