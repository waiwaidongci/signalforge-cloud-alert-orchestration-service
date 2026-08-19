package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/acme/signalforge/internal/incident/domain"
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

func (r *SQLRepository) Create(ctx context.Context, incident domain.Incident) error {
	_, err := r.store.Exec(ctx, `
		INSERT INTO incidents (
			id, fingerprint, title, severity, status, source_id, alert_count, first_seen_at, last_seen_at,
			acknowledged_at, closed_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		incident.ID, incident.Fingerprint, incident.Title, incident.Severity.String(), string(incident.Status), incident.SourceID, incident.AlertCount,
		db.NowString(incident.FirstSeenAt), db.NowString(incident.LastSeenAt), emptyTime(incident.AcknowledgedAt), emptyTime(incident.ClosedAt),
		db.NowString(incident.CreatedAt), db.NowString(incident.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("create incident: %w", err)
	}
	return nil
}

func (r *SQLRepository) Update(ctx context.Context, incident domain.Incident) error {
	result, err := r.store.Exec(ctx, `
		UPDATE incidents SET title = ?, severity = ?, status = ?, source_id = ?, alert_count = ?,
			first_seen_at = ?, last_seen_at = ?, acknowledged_at = ?, closed_at = ?, updated_at = ?
		WHERE id = ?`,
		incident.Title, incident.Severity.String(), string(incident.Status), incident.SourceID, incident.AlertCount,
		db.NowString(incident.FirstSeenAt), db.NowString(incident.LastSeenAt), emptyTime(incident.AcknowledgedAt), emptyTime(incident.ClosedAt),
		db.NowString(incident.UpdatedAt), incident.ID,
	)
	return store.NotModified(result, err, "INCIDENT_NOT_FOUND", "事件组不存在")
}

func (r *SQLRepository) FindByID(ctx context.Context, id string) (domain.Incident, error) {
	return scanIncident(r.store.QueryRow(ctx, incidentSelect+` WHERE id = ?`, id))
}

func (r *SQLRepository) FindByFingerprint(ctx context.Context, fingerprint string) (domain.Incident, error) {
	return scanIncident(r.store.QueryRow(ctx, incidentSelect+` WHERE fingerprint = ? ORDER BY created_at DESC LIMIT 1`, fingerprint))
}

func (r *SQLRepository) List(ctx context.Context, filter domain.ListFilter, limit, offset int) ([]domain.Incident, int, error) {
	where, args := buildIncidentWhere(filter)
	count, err := r.count(ctx, where, args...)
	if err != nil {
		return nil, 0, err
	}
	query := incidentSelect + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.store.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var incidents []domain.Incident
	for rows.Next() {
		incident, err := scanIncident(rows)
		if err != nil {
			return nil, 0, err
		}
		incidents = append(incidents, incident)
	}
	return incidents, count, rows.Err()
}

func (r *SQLRepository) IncrementAlertCount(ctx context.Context, id string, lastSeenAt time.Time) error {
	result, err := r.store.Exec(ctx, `UPDATE incidents SET alert_count = alert_count + 1, last_seen_at = ?, updated_at = ? WHERE id = ?`,
		db.NowString(lastSeenAt), db.NowString(time.Now().UTC()), id)
	return store.NotModified(result, err, "INCIDENT_NOT_FOUND", "事件组不存在")
}

const incidentSelect = `
	SELECT id, fingerprint, title, severity, status, source_id, alert_count, first_seen_at, last_seen_at,
		acknowledged_at, closed_at, created_at, updated_at
	FROM incidents`

func buildIncidentWhere(filter domain.ListFilter) (string, []any) {
	where := " WHERE 1=1"
	args := []any{}
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
	if filter.Query != "" {
		where += " AND title LIKE ?"
		args = append(args, "%"+filter.Query+"%")
	}
	return where, args
}

func (r *SQLRepository) count(ctx context.Context, where string, args ...any) (int, error) {
	var total int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM incidents`+where, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count incidents: %w", err)
	}
	return total, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanIncident(row scanner) (domain.Incident, error) {
	var incident domain.Incident
	var severityValue, status, firstSeenAt, lastSeenAt, ackAt, closedAt, createdAt, updatedAt string
	if err := row.Scan(
		&incident.ID, &incident.Fingerprint, &incident.Title, &severityValue, &status, &incident.SourceID, &incident.AlertCount,
		&firstSeenAt, &lastSeenAt, &ackAt, &closedAt, &createdAt, &updatedAt,
	); err != nil {
		return domain.Incident{}, normalizeNotFound(err)
	}
	incident.Severity = severity.Severity(severityValue)
	incident.Status = domain.Status(status)
	incident.FirstSeenAt, _ = db.ParseTime(firstSeenAt)
	incident.LastSeenAt, _ = db.ParseTime(lastSeenAt)
	incident.AcknowledgedAt = parseOptionalTime(ackAt)
	incident.ClosedAt = parseOptionalTime(closedAt)
	incident.CreatedAt, _ = db.ParseTime(createdAt)
	incident.UpdatedAt, _ = db.ParseTime(updatedAt)
	return incident, nil
}

func parseOptionalTime(raw string) time.Time {
	value, _ := db.ParseTime(raw)
	return value
}

func emptyTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return db.NowString(value)
}

func normalizeNotFound(err error) error {
	if err == sql.ErrNoRows {
		return apperr.NotFound("INCIDENT_NOT_FOUND", "事件组不存在")
	}
	return err
}
