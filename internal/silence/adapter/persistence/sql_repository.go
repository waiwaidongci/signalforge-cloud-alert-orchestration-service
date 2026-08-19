package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/db"
	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/shared/severity"
	"github.com/acme/signalforge/internal/shared/store"
	"github.com/acme/signalforge/internal/silence/domain"
)

type SQLRepository struct {
	store *store.Store
}

func NewSQLRepository(raw *sql.DB) *SQLRepository {
	return &SQLRepository{store: store.New(raw)}
}

func (r *SQLRepository) CreateSilence(ctx context.Context, silence domain.Silence) error {
	matcherJSON, _ := json.Marshal(silence.Matcher)
	_, err := r.store.Exec(ctx, `
		INSERT INTO silences (id, name, description, matcher_json, starts_at, ends_at, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		silence.ID, silence.Name, silence.Description, string(matcherJSON), db.NowString(silence.StartsAt), db.NowString(silence.EndsAt), silence.CreatedBy, db.NowString(silence.CreatedAt), db.NowString(silence.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("create silence: %w", err)
	}
	return nil
}

func (r *SQLRepository) UpdateSilence(ctx context.Context, silence domain.Silence) error {
	matcherJSON, _ := json.Marshal(silence.Matcher)
	result, err := r.store.Exec(ctx, `
		UPDATE silences SET name = ?, description = ?, matcher_json = ?, starts_at = ?, ends_at = ?, created_by = ?, updated_at = ?
		WHERE id = ?`,
		silence.Name, silence.Description, string(matcherJSON), db.NowString(silence.StartsAt), db.NowString(silence.EndsAt), silence.CreatedBy, db.NowString(silence.UpdatedAt), silence.ID,
	)
	return store.NotModified(result, err, "SILENCE_NOT_FOUND", "静默规则不存在")
}

func (r *SQLRepository) FindSilenceByID(ctx context.Context, id string) (domain.Silence, error) {
	return scanSilence(r.store.QueryRow(ctx, silenceSelect+` WHERE id = ?`, id))
}

func (r *SQLRepository) ListSilences(ctx context.Context, activeOnly bool, limit, offset int) ([]domain.Silence, int, error) {
	where := " WHERE 1=1"
	args := []any{}
	if activeOnly {
		where += " AND ends_at >= ?"
		args = append(args, db.NowString(timeNow()))
	}
	count, err := r.countSilences(ctx, where, args...)
	if err != nil {
		return nil, 0, err
	}
	query := silenceSelect + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.store.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var silences []domain.Silence
	for rows.Next() {
		silence, err := scanSilence(rows)
		if err != nil {
			return nil, 0, err
		}
		silences = append(silences, silence)
	}
	return silences, count, rows.Err()
}

func (r *SQLRepository) DeleteSilence(ctx context.Context, id string) error {
	result, err := r.store.Exec(ctx, `DELETE FROM silences WHERE id = ?`, id)
	return store.NotModified(result, err, "SILENCE_NOT_FOUND", "静默规则不存在")
}

func (r *SQLRepository) ActiveSilences(ctx context.Context, at time.Time) ([]domain.Silence, error) {
	rows, err := r.store.Query(ctx, silenceSelect+` WHERE starts_at <= ? AND ends_at >= ? ORDER BY starts_at ASC`, db.NowString(at), db.NowString(at))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var silences []domain.Silence
	for rows.Next() {
		silence, err := scanSilence(rows)
		if err != nil {
			return nil, err
		}
		silences = append(silences, silence)
	}
	return silences, rows.Err()
}

func (r *SQLRepository) CreateSuppression(ctx context.Context, rule domain.SuppressionRule) error {
	matcherJSON, _ := json.Marshal(rule.SourceMatcher)
	_, err := r.store.Exec(ctx, `
		INSERT INTO suppression_rules (id, name, description, source_matcher_json, high_severity, low_severity, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rule.ID, rule.Name, rule.Description, string(matcherJSON), rule.HighSeverity.String(), rule.LowSeverity.String(), boolInt(rule.Enabled), db.NowString(rule.CreatedAt), db.NowString(rule.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("create suppression: %w", err)
	}
	return nil
}

func (r *SQLRepository) UpdateSuppression(ctx context.Context, rule domain.SuppressionRule) error {
	matcherJSON, _ := json.Marshal(rule.SourceMatcher)
	result, err := r.store.Exec(ctx, `
		UPDATE suppression_rules SET name = ?, description = ?, source_matcher_json = ?, high_severity = ?, low_severity = ?, enabled = ?, updated_at = ?
		WHERE id = ?`,
		rule.Name, rule.Description, string(matcherJSON), rule.HighSeverity.String(), rule.LowSeverity.String(), boolInt(rule.Enabled), db.NowString(rule.UpdatedAt), rule.ID,
	)
	return store.NotModified(result, err, "SUPPRESSION_NOT_FOUND", "抑制规则不存在")
}

func (r *SQLRepository) FindSuppressionByID(ctx context.Context, id string) (domain.SuppressionRule, error) {
	return scanSuppression(r.store.QueryRow(ctx, suppressionSelect+` WHERE id = ?`, id))
}

func (r *SQLRepository) ListSuppressions(ctx context.Context, enabledOnly bool, limit, offset int) ([]domain.SuppressionRule, int, error) {
	where := " WHERE 1=1"
	args := []any{}
	if enabledOnly {
		where += " AND enabled = ?"
		args = append(args, 1)
	}
	count, err := r.countSuppressions(ctx, where, args...)
	if err != nil {
		return nil, 0, err
	}
	query := suppressionSelect + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.store.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var rules []domain.SuppressionRule
	for rows.Next() {
		rule, err := scanSuppression(rows)
		if err != nil {
			return nil, 0, err
		}
		rules = append(rules, rule)
	}
	return rules, count, rows.Err()
}

func (r *SQLRepository) DeleteSuppression(ctx context.Context, id string) error {
	result, err := r.store.Exec(ctx, `DELETE FROM suppression_rules WHERE id = ?`, id)
	return store.NotModified(result, err, "SUPPRESSION_NOT_FOUND", "抑制规则不存在")
}

func (r *SQLRepository) ActiveSuppressions(ctx context.Context) ([]domain.SuppressionRule, error) {
	rows, err := r.store.Query(ctx, suppressionSelect+` WHERE enabled = 1 ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []domain.SuppressionRule
	for rows.Next() {
		rule, err := scanSuppression(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

const silenceSelect = `
	SELECT id, name, description, matcher_json, starts_at, ends_at, created_by, created_at, updated_at
	FROM silences`

const suppressionSelect = `
	SELECT id, name, description, source_matcher_json, high_severity, low_severity, enabled, created_at, updated_at
	FROM suppression_rules`

type scanner interface {
	Scan(dest ...any) error
}

func scanSilence(row scanner) (domain.Silence, error) {
	var silence domain.Silence
	var matcherJSON, startsAt, endsAt, createdAt, updatedAt string
	if err := row.Scan(&silence.ID, &silence.Name, &silence.Description, &matcherJSON, &startsAt, &endsAt, &silence.CreatedBy, &createdAt, &updatedAt); err != nil {
		return domain.Silence{}, normalizeNotFound(err)
	}
	silence.Matcher, _ = matcher.DecodeSelector(matcherJSON)
	silence.StartsAt, _ = db.ParseTime(startsAt)
	silence.EndsAt, _ = db.ParseTime(endsAt)
	silence.CreatedAt, _ = db.ParseTime(createdAt)
	silence.UpdatedAt, _ = db.ParseTime(updatedAt)
	return silence, nil
}

func scanSuppression(row scanner) (domain.SuppressionRule, error) {
	var rule domain.SuppressionRule
	var matcherJSON, high, low, createdAt, updatedAt string
	var enabled int
	if err := row.Scan(&rule.ID, &rule.Name, &rule.Description, &matcherJSON, &high, &low, &enabled, &createdAt, &updatedAt); err != nil {
		return domain.SuppressionRule{}, normalizeNotFound(err)
	}
	rule.SourceMatcher, _ = matcher.DecodeSelector(matcherJSON)
	rule.HighSeverity = severity.Severity(high)
	rule.LowSeverity = severity.Severity(low)
	rule.Enabled = enabled != 0
	rule.CreatedAt, _ = db.ParseTime(createdAt)
	rule.UpdatedAt, _ = db.ParseTime(updatedAt)
	return rule, nil
}

func (r *SQLRepository) countSilences(ctx context.Context, where string, args ...any) (int, error) {
	var total int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM silences`+where, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count silences: %w", err)
	}
	return total, nil
}

func (r *SQLRepository) countSuppressions(ctx context.Context, where string, args ...any) (int, error) {
	var total int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM suppression_rules`+where, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count suppressions: %w", err)
	}
	return total, nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func timeNow() time.Time {
	return time.Now().UTC()
}

func normalizeNotFound(err error) error {
	if err == sql.ErrNoRows {
		return apperr.NotFound("NOT_FOUND", "资源不存在")
	}
	return err
}
