package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/acme/signalforge/internal/shared/db"
)

type Store struct {
	DB     *sql.DB
	Driver string
}

func New(raw *sql.DB) *Store {
	return &Store{DB: raw, Driver: db.DriverName(raw)}
}

func (s *Store) Rebind(query string) string {
	return db.Rebind(s.Driver, query)
}

func (s *Store) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	result, err := s.DB.ExecContext(ctx, s.Rebind(query), args...)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}
	return result, nil
}

func (s *Store) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	rows, err := s.DB.QueryContext(ctx, s.Rebind(query), args...)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}
	return rows, nil
}

func (s *Store) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return s.DB.QueryRowContext(ctx, s.Rebind(query), args...)
}

func EncodeMap(value map[string]string) string {
	raw, _ := json.Marshal(value)
	return string(raw)
}

func DecodeMap(raw string) (map[string]string, error) {
	if raw == "" {
		return map[string]string{}, nil
	}
	var value map[string]string
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return nil, fmt.Errorf("decode map: %w", err)
	}
	return value, nil
}

func DecodePayload(raw string) (map[string]any, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, ErrEmptyPayload
	}
	if raw == "null" {
		return nil, ErrEmptyPayload
	}
	var value map[string]any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}
	if value == nil {
		return nil, ErrEmptyPayload
	}
	return value, nil
}

func EncodeAny(value any) string {
	raw, _ := json.Marshal(value)
	return string(raw)
}

func DecodeAny(raw string, target any) error {
	if raw == "" {
		return nil
	}
	return json.Unmarshal([]byte(raw), target)
}
