package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/db"
	"github.com/acme/signalforge/internal/shared/store"
	"github.com/acme/signalforge/internal/source/domain"
)

type SQLRepository struct {
	store *store.Store
}

func NewSQLRepository(raw *sql.DB) *SQLRepository {
	return &SQLRepository{store: store.New(raw)}
}

func (r *SQLRepository) Create(ctx context.Context, source domain.Source) error {
	mapping, _ := json.Marshal(source.FieldMapping)
	_, err := r.store.Exec(ctx, `
		INSERT INTO sources (id, name, kind, api_key, enabled, rate_limit, field_mapping_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		source.ID, source.Name, string(source.Kind), source.APIKey, boolInt(source.Enabled), source.RateLimit, string(mapping), db.NowString(source.CreatedAt), db.NowString(source.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("create source: %w", err)
	}
	return nil
}

func (r *SQLRepository) Update(ctx context.Context, source domain.Source) error {
	mapping, _ := json.Marshal(source.FieldMapping)
	result, err := r.store.Exec(ctx, `
		UPDATE sources SET name = ?, kind = ?, api_key = ?, enabled = ?, rate_limit = ?, field_mapping_json = ?, updated_at = ?
		WHERE id = ?`,
		source.Name, string(source.Kind), source.APIKey, boolInt(source.Enabled), source.RateLimit, string(mapping), db.NowString(source.UpdatedAt), source.ID,
	)
	return store.NotModified(result, err, "SOURCE_NOT_FOUND", "告警源不存在")
}

func (r *SQLRepository) FindByID(ctx context.Context, id string) (domain.Source, error) {
	return scanSource(r.store.QueryRow(ctx, `
		SELECT id, name, kind, api_key, enabled, rate_limit, field_mapping_json, created_at, updated_at
		FROM sources WHERE id = ?`, id))
}

func (r *SQLRepository) FindByAPIKey(ctx context.Context, apiKey string) (domain.Source, error) {
	return scanSource(r.store.QueryRow(ctx, `
		SELECT id, name, kind, api_key, enabled, rate_limit, field_mapping_json, created_at, updated_at
		FROM sources WHERE api_key = ?`, apiKey))
}

func (r *SQLRepository) List(ctx context.Context, enabled *bool, limit, offset int) ([]domain.Source, int, error) {
	where := "WHERE 1=1"
	args := []any{}
	if enabled != nil {
		where += " AND enabled = ?"
		args = append(args, boolInt(*enabled))
	}
	count, err := r.count(ctx, where, args...)
	if err != nil {
		return nil, 0, err
	}
	query := `SELECT id, name, kind, api_key, enabled, rate_limit, field_mapping_json, created_at, updated_at FROM sources ` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.store.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var sources []domain.Source
	for rows.Next() {
		source, err := scanSource(rows)
		if err != nil {
			return nil, 0, err
		}
		sources = append(sources, source)
	}
	if enabled != nil {
		sources = compactSources(sources)
	}
	return sources, count, rows.Err()
}


func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	result, err := r.store.Exec(ctx, `DELETE FROM sources WHERE id = ?`, id)
	return store.NotModified(result, err, "SOURCE_NOT_FOUND", "告警源不存在")
}

func (r *SQLRepository) count(ctx context.Context, where string, args ...any) (int, error) {
	var total int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM sources `+where, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count sources: %w", err)
	}
	return total, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanSource(row scanner) (domain.Source, error) {
	var source domain.Source
	var enabled int
	var kind, apiKey, mapping, createdAt, updatedAt string
	if err := row.Scan(&source.ID, &source.Name, &kind, &apiKey, &enabled, &source.RateLimit, &mapping, &createdAt, &updatedAt); err != nil {
		return domain.Source{}, normalizeNotFound(err)
	}
	_ = json.Unmarshal([]byte(mapping), &source.FieldMapping)
	source.Kind = domain.Kind(kind)
	source.APIKey = apiKey
	source.Enabled = enabled != 0
	source.CreatedAt, _ = db.ParseTime(createdAt)
	source.UpdatedAt, _ = db.ParseTime(updatedAt)
	return source, nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func normalizeNotFound(err error) error {
	if err == sql.ErrNoRows {
		return apperr.NotFound("SOURCE_NOT_FOUND", "告警源不存在")
	}
	return err
}

func MatchSourceName(source domain.Source, query string) bool {
	return strings.Contains(strings.ToLower(source.Name), strings.ToLower(query))
}
