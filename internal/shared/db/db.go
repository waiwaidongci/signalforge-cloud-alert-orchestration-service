package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/acme/signalforge/internal/shared/config"
	_ "modernc.org/sqlite"
)

func Open(cfg config.DatabaseConfig) (*sql.DB, error) {
	if cfg.Driver == "" {
		cfg.Driver = "sqlite"
	}
	if cfg.Driver == "sqlite" {
		if cfg.DSN == "" || cfg.DSN == ":memory:" {
			return openSQLite(cfg)
		}
		if err := os.MkdirAll(filepath.Dir(cfg.DSN), 0o755); err != nil {
			return nil, fmt.Errorf("create sqlite directory: %w", err)
		}
	}
	raw, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if cfg.MaxOpenConns > 0 {
		raw.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns >= 0 {
		raw.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime.Value() > 0 {
		raw.SetConnMaxLifetime(cfg.ConnMaxLifetime.Value())
	}
	if cfg.Driver == "sqlite" {
		raw.SetMaxOpenConns(1)
		raw.SetMaxIdleConns(1)
		raw.SetConnMaxLifetime(0)
	}
	if err := raw.Ping(); err != nil {
		_ = raw.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return raw, nil
}

func openSQLite(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.DSN
	if dsn == "" {
		dsn = "file:signalforge.db?mode=memory&cache=shared"
	}
	raw, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	raw.SetMaxOpenConns(1)
	raw.SetMaxIdleConns(1)
	raw.SetConnMaxLifetime(0)
	if err := raw.Ping(); err != nil {
		_ = raw.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	return raw, nil
}

func NowString(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func ParseTime(raw string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, raw)
}
