package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/acme/signalforge/internal/escalation/domain"
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

func (r *SQLRepository) Create(ctx context.Context, policy domain.Policy) error {
	matchJSON, _ := json.Marshal(policy.Matcher)
	routesJSON, _ := json.Marshal(policy.Routes)
	_, err := r.store.Exec(ctx, `
		INSERT INTO escalation_policies (id, name, description, match_json, wait_seconds, repeat_seconds, max_repeats, routes_json, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		policy.ID, policy.Name, policy.Description, string(matchJSON), policy.WaitSeconds, policy.RepeatSeconds, policy.MaxRepeats, string(routesJSON), boolInt(policy.Enabled), db.NowString(policy.CreatedAt), db.NowString(policy.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("create escalation policy: %w", err)
	}
	return nil
}

func (r *SQLRepository) Update(ctx context.Context, policy domain.Policy) error {
	matchJSON, _ := json.Marshal(policy.Matcher)
	routesJSON, _ := json.Marshal(policy.Routes)
	result, err := r.store.Exec(ctx, `
		UPDATE escalation_policies SET name = ?, description = ?, match_json = ?, wait_seconds = ?, repeat_seconds = ?, max_repeats = ?, routes_json = ?, enabled = ?, updated_at = ?
		WHERE id = ?`,
		policy.Name, policy.Description, string(matchJSON), policy.WaitSeconds, policy.RepeatSeconds, policy.MaxRepeats, string(routesJSON), boolInt(policy.Enabled), db.NowString(policy.UpdatedAt), policy.ID,
	)
	return store.NotModified(result, err, "ESCALATION_POLICY_NOT_FOUND", "升级策略不存在")
}

func (r *SQLRepository) FindByID(ctx context.Context, id string) (domain.Policy, error) {
	return scanPolicy(r.store.QueryRow(ctx, policySelect+` WHERE id = ?`, id))
}

func (r *SQLRepository) List(ctx context.Context, enabledOnly bool, limit, offset int) ([]domain.Policy, int, error) {
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
	query := policySelect + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.store.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var policies []domain.Policy
	for rows.Next() {
		policy, err := scanPolicy(rows)
		if err != nil {
			return nil, 0, err
		}
		policies = append(policies, policy)
	}
	return policies, count, rows.Err()
}

func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	result, err := r.store.Exec(ctx, `DELETE FROM escalation_policies WHERE id = ?`, id)
	return store.NotModified(result, err, "ESCALATION_POLICY_NOT_FOUND", "升级策略不存在")
}

func (r *SQLRepository) Active(ctx context.Context) ([]domain.Policy, error) {
	rows, err := r.store.Query(ctx, policySelect+` WHERE enabled = 1 ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var policies []domain.Policy
	for rows.Next() {
		policy, err := scanPolicy(rows)
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}
	return policies, rows.Err()
}

const policySelect = `
	SELECT id, name, description, match_json, wait_seconds, repeat_seconds, max_repeats, routes_json, enabled, created_at, updated_at
	FROM escalation_policies`

type scanner interface {
	Scan(dest ...any) error
}

func scanPolicy(row scanner) (domain.Policy, error) {
	var policy domain.Policy
	var matchJSON, routesJSON, createdAt, updatedAt string
	var enabled int
	if err := row.Scan(&policy.ID, &policy.Name, &policy.Description, &matchJSON, &policy.WaitSeconds, &policy.RepeatSeconds, &policy.MaxRepeats, &routesJSON, &enabled, &createdAt, &updatedAt); err != nil {
		return domain.Policy{}, normalizeNotFound(err)
	}
	policy.Matcher, _ = matcher.DecodeSelector(matchJSON)
	_ = json.Unmarshal([]byte(routesJSON), &policy.Routes)
	policy.Enabled = enabled != 0
	policy.CreatedAt, _ = db.ParseTime(createdAt)
	policy.UpdatedAt, _ = db.ParseTime(updatedAt)
	return policy, nil
}

func (r *SQLRepository) count(ctx context.Context, where string, args ...any) (int, error) {
	var total int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM escalation_policies`+where, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count escalation policies: %w", err)
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
		return apperr.NotFound("ESCALATION_POLICY_NOT_FOUND", "升级策略不存在")
	}
	return err
}
