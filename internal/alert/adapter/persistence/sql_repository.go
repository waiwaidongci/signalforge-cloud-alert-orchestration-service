package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/acme/signalforge/internal/alert/domain"
	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/db"
	"github.com/acme/signalforge/internal/shared/severity"
	"github.com/acme/signalforge/internal/shared/store"
)

type SQLRepository struct {
	store *store.Store
}

func NewSQLRepository(raw *sql.DB) *SQLRepository {
	return &SQLRepository{store: store.New(raw)}
}

func (r *SQLRepository) Create(ctx context.Context, alert domain.Alert) error {
	labels, _ := json.Marshal(alert.Labels)
	annotations, _ := json.Marshal(alert.Annotations)
	_, err := r.store.Exec(ctx, `
		INSERT INTO alerts (
			id, source_id, external_id, fingerprint, resource, severity, status, title, description,
			labels_json, annotations_json, received_at, last_occurred_at, dedup_count, incident_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		alert.ID, alert.SourceID, alert.ExternalID, alert.Fingerprint, alert.Resource, alert.Severity.String(), string(alert.Status), alert.Title, alert.Description,
		string(labels), string(annotations), db.NowString(alert.ReceivedAt), db.NowString(alert.LastOccurredAt), alert.DedupCount, alert.IncidentID, db.NowString(alert.CreatedAt), db.NowString(alert.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("create alert: %w", err)
	}
	return nil
}

func (r *SQLRepository) Update(ctx context.Context, alert domain.Alert) error {
	labels, _ := json.Marshal(alert.Labels)
	annotations, _ := json.Marshal(alert.Annotations)
	result, err := r.store.Exec(ctx, `
		UPDATE alerts SET fingerprint = ?, resource = ?, severity = ?, status = ?, title = ?, description = ?,
			labels_json = ?, annotations_json = ?, received_at = ?, last_occurred_at = ?, dedup_count = ?, incident_id = ?, updated_at = ?
		WHERE id = ?`,
		alert.Fingerprint, alert.Resource, alert.Severity.String(), string(alert.Status), alert.Title, alert.Description,
		string(labels), string(annotations), db.NowString(alert.ReceivedAt), db.NowString(alert.LastOccurredAt), alert.DedupCount, alert.IncidentID, db.NowString(alert.UpdatedAt), alert.ID,
	)
	return store.NotModified(result, err, "ALERT_NOT_FOUND", "告警不存在")
}

func (r *SQLRepository) FindByID(ctx context.Context, id string) (domain.Alert, error) {
	return scanAlert(r.store.QueryRow(ctx, alertSelect+` WHERE id = ?`, id))
}

func (r *SQLRepository) FindBySourceExternalID(ctx context.Context, sourceID, externalID string) (domain.Alert, error) {
	return scanAlert(r.store.QueryRow(ctx, alertSelect+` WHERE source_id = ? AND external_id = ?`, sourceID, externalID))
}

func (r *SQLRepository) FindByFingerprint(ctx context.Context, fingerprint string, limit int) ([]domain.Alert, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.store.Query(ctx, alertSelect+` WHERE fingerprint = ? ORDER BY last_occurred_at DESC LIMIT ?`, fingerprint, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var alerts []domain.Alert
	for rows.Next() {
		alert, err := scanAlert(rows)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, rows.Err()
}

func (r *SQLRepository) List(ctx context.Context, filter domain.ListFilter, limit, offset int) ([]domain.Alert, int, error) {
	where, args := buildAlertWhere(filter)
	count, err := r.count(ctx, where, args...)
	if err != nil {
		return nil, 0, err
	}
	query := alertSelect + where + ` ORDER BY last_occurred_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.store.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var alerts []domain.Alert
	for rows.Next() {
		alert, err := scanAlert(rows)
		if err != nil {
			return nil, 0, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, count, rows.Err()
}

func (r *SQLRepository) UpdateIncidentID(ctx context.Context, ids []string, incidentID string) error {
	if len(ids) == 0 {
		return nil
	}
	args := []any{incidentID, db.NowString(timeNow())}
	for _, id := range ids {
		args = append(args, id)
	}
	query := `UPDATE alerts SET incident_id = ?, updated_at = ? WHERE id IN (` + placeholders(len(ids)) + `)`
	result, err := r.store.Exec(ctx, query, args...)
	return store.NotModified(result, err, "ALERT_NOT_FOUND", "告警不存在")
}

func (r *SQLRepository) BatchUpdateStatus(ctx context.Context, ids []string, status domain.Status, at time.Time) error {
	if len(ids) == 0 {
		return nil
	}
	args := []any{string(normalizeStatus(status)), db.NowString(at)}
	for _, id := range ids {
		args = append(args, id)
	}
	query := `UPDATE alerts SET status = ?, updated_at = ? WHERE id IN (` + placeholders(len(ids)) + `)`
	result, err := r.store.Exec(ctx, query, args...)
	return store.NotModified(result, err, "ALERT_NOT_FOUND", "告警不存在")
}

func normalizeStatus(status domain.Status) domain.Status {
	switch status {
	case domain.StatusAcknowledged, domain.StatusFiring, domain.StatusResolved, domain.StatusSuppressed:
		return status
	default:
		return domain.StatusFiring
	}
}

func (r *SQLRepository) CountByIncident(ctx context.Context, incidentID string, status domain.Status) (int, error) {
	var count int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM alerts WHERE incident_id = ? AND status = ?`, incidentID, string(status)).Scan(&count); err != nil {
		return 0, fmt.Errorf("count alerts by incident: %w", err)
	}
	return count, nil
}

func (r *SQLRepository) HighestActiveSeverity(ctx context.Context) (severity.Severity, error) {
	rows, err := r.store.Query(ctx, `SELECT severity FROM alerts WHERE status = ?`, string(domain.StatusFiring))
	if err != nil {
		return severity.Info, err
	}
	defer rows.Close()
	highest := severity.Info
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return severity.Info, err
		}
		current := severity.Severity(raw)
		if severity.HigherOrEqual(current, highest) {
			highest = current
		}
	}
	return highest, rows.Err()
}

const alertSelect = `
	SELECT id, source_id, external_id, fingerprint, resource, severity, status, title, description,
		labels_json, annotations_json, received_at, last_occurred_at, dedup_count, incident_id, created_at, updated_at
	FROM alerts`

func buildAlertWhere(filter domain.ListFilter) (string, []any) {
	where := " WHERE 1=1"
	args := []any{}
	if filter.SourceID != "" {
		where += " AND source_id = ?"
		args = append(args, filter.SourceID)
	}
	if filter.Fingerprint != "" {
		where += " AND fingerprint = ?"
		args = append(args, filter.Fingerprint)
	}
	if filter.Status != "" {
		where += " AND status = ?"
		args = append(args, string(filter.Status))
	}
	if filter.Severity != "" {
		where += " AND severity = ?"
		args = append(args, filter.Severity)
	}
	if filter.IncidentID != "" {
		where += " AND incident_id = ?"
		args = append(args, filter.IncidentID)
	}
	if filter.Query != "" {
		where += " AND (title LIKE ? OR description LIKE ? OR resource LIKE ?)"
		like := "%" + filter.Query + "%"
		args = append(args, like, like, like)
	}
	if !filter.From.IsZero() {
		where += " AND received_at >= ?"
		args = append(args, db.NowString(filter.From))
	}
	if !filter.To.IsZero() {
		where += " AND received_at <= ?"
		args = append(args, db.NowString(filter.To))
	}
	return where, args
}

func (r *SQLRepository) count(ctx context.Context, where string, args ...any) (int, error) {
	var total int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM alerts`+where, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count alerts: %w", err)
	}
	return total, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanAlert(row scanner) (domain.Alert, error) {
	var alert domain.Alert
	var severity, status, labels, annotations, receivedAt, lastOccurredAt, createdAt, updatedAt string
	if err := row.Scan(
		&alert.ID, &alert.SourceID, &alert.ExternalID, &alert.Fingerprint, &alert.Resource, &severity, &status, &alert.Title, &alert.Description,
		&labels, &annotations, &receivedAt, &lastOccurredAt, &alert.DedupCount, &alert.IncidentID, &createdAt, &updatedAt,
	); err != nil {
		return domain.Alert{}, normalizeNotFound(err)
	}
	_ = json.Unmarshal([]byte(labels), &alert.Labels)
	_ = json.Unmarshal([]byte(annotations), &alert.Annotations)
	alert.Severity = severityValue(severity)
	alert.Status = domain.Status(status)
	alert.ReceivedAt, _ = db.ParseTime(receivedAt)
	alert.LastOccurredAt, _ = db.ParseTime(lastOccurredAt)
	alert.CreatedAt, _ = db.ParseTime(createdAt)
	alert.UpdatedAt, _ = db.ParseTime(updatedAt)
	return alert, nil
}

func severityValue(raw string) severity.Severity {
	return severity.Severity(raw)
}

func placeholders(count int) string {
	values := make([]string, count)
	for i := range values {
		values[i] = "?"
	}
	return strings.Join(values, ",")
}

func normalizeNotFound(err error) error {
	if err == sql.ErrNoRows {
		return apperr.NotFound("ALERT_NOT_FOUND", "告警不存在")
	}
	return err
}

func timeNow() time.Time {
	return time.Now().UTC()
}
