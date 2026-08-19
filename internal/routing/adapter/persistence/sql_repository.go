package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/acme/signalforge/internal/routing/domain"
	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/db"
	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/shared/store"
)

type SQLRepository struct {
	store *store.Store
}

func NewSQLRepository(raw *sql.DB) *SQLRepository {
	return &SQLRepository{store: store.New(raw)}
}

func (r *SQLRepository) Create(ctx context.Context, rule domain.Rule) error {
	matchJSON, _ := json.Marshal(rule.Match)
	channelsJSON, _ := json.Marshal(rule.Channels)
	_, err := r.store.Exec(ctx, `
		INSERT INTO routing_rules (id, name, priority, match_json, channels_json, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		rule.ID, rule.Name, rule.Priority, string(matchJSON), string(channelsJSON), boolInt(rule.Enabled), db.NowString(rule.CreatedAt), db.NowString(rule.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("create routing rule: %w", err)
	}
	return nil
}

func (r *SQLRepository) Update(ctx context.Context, rule domain.Rule) error {
	matchJSON, _ := json.Marshal(rule.Match)
	channelsJSON, _ := json.Marshal(rule.Channels)
	result, err := r.store.Exec(ctx, `
		UPDATE routing_rules SET name = ?, priority = ?, match_json = ?, channels_json = ?, enabled = ?, updated_at = ?
		WHERE id = ?`,
		rule.Name, rule.Priority, string(matchJSON), string(channelsJSON), boolInt(rule.Enabled), db.NowString(rule.UpdatedAt), rule.ID,
	)
	return store.NotModified(result, err, "ROUTING_RULE_NOT_FOUND", "路由策略不存在")
}

func (r *SQLRepository) FindByID(ctx context.Context, id string) (domain.Rule, error) {
	return scanRule(r.store.QueryRow(ctx, ruleSelect+` WHERE id = ?`, id))
}

func (r *SQLRepository) List(ctx context.Context, enabledOnly bool, limit, offset int) ([]domain.Rule, int, error) {
	where := " WHERE 1=1"
	args := []any{}
	if enabledOnly {
		where += " AND enabled = ?"
		args = append(args, 1)
	}
	count, err := r.count(ctx, where, args...)
	if err != nil {
		return nil, 0, err
	}
	query := ruleSelect + where + ` ORDER BY priority DESC, created_at ASC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.store.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var rules []domain.Rule
	for rows.Next() {
		rule, err := scanRule(rows)
		if err != nil {
			return nil, 0, err
		}
		rules = append(rules, rule)
	}
	return rules, count, rows.Err()
}

func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	result, err := r.store.Exec(ctx, `DELETE FROM routing_rules WHERE id = ?`, id)
	return store.NotModified(result, err, "ROUTING_RULE_NOT_FOUND", "路由策略不存在")
}

func (r *SQLRepository) Active(ctx context.Context) ([]domain.Rule, error) {
	rows, err := r.store.Query(ctx, ruleSelect+` WHERE enabled = 1 ORDER BY priority DESC, created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []domain.Rule
	for rows.Next() {
		rule, err := scanRule(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

const ruleSelect = `
	SELECT id, name, priority, match_json, channels_json, enabled, created_at, updated_at
	FROM routing_rules`

type scanner interface {
	Scan(dest ...any) error
}

func scanRule(row scanner) (domain.Rule, error) {
	var rule domain.Rule
	var matchJSON, channelsJSON, createdAt, updatedAt string
	var enabled int
	if err := row.Scan(&rule.ID, &rule.Name, &rule.Priority, &matchJSON, &channelsJSON, &enabled, &createdAt, &updatedAt); err != nil {
		return domain.Rule{}, normalizeNotFound(err)
	}
	rule.Match, _ = matcher.DecodeSelector(matchJSON)
	_ = json.Unmarshal([]byte(channelsJSON), &rule.Channels)
	rule.Enabled = enabled != 0
	rule.CreatedAt, _ = db.ParseTime(createdAt)
	rule.UpdatedAt, _ = db.ParseTime(updatedAt)
	return rule, nil
}

func (r *SQLRepository) count(ctx context.Context, where string, args ...any) (int, error) {
	var total int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM routing_rules`+where, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count routing rules: %w", err)
	}
	return total, nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func normalizeNotFound(err error) error {
	if err == sql.ErrNoRows {
		return apperr.NotFound("ROUTING_RULE_NOT_FOUND", "路由策略不存在")
	}
	return err
}
