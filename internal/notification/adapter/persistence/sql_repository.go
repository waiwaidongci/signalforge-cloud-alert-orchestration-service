package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/acme/signalforge/internal/notification/domain"
	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/db"
	"github.com/acme/signalforge/internal/shared/store"
)

type SQLRepository struct {
	store *store.Store
}

func NewSQLRepository(raw *sql.DB) *SQLRepository {
	return &SQLRepository{store: store.New(raw)}
}

func (r *SQLRepository) Create(ctx context.Context, notification domain.Notification) error {
	payload, _ := json.Marshal(notification.Payload)
	_, err := r.store.Exec(ctx, `
		INSERT INTO notifications (id, incident_id, alert_id, channel, destination, status, payload_json, error_message, sent_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		notification.ID, notification.IncidentID, notification.AlertID, notification.Channel, notification.Destination, string(notification.Status), string(payload), notification.ErrorMessage, emptyTime(notification.SentAt), db.NowString(notification.CreatedAt), db.NowString(notification.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *SQLRepository) Update(ctx context.Context, notification domain.Notification) error {
	payload, _ := json.Marshal(notification.Payload)
	result, err := r.store.Exec(ctx, `
		UPDATE notifications SET channel = ?, destination = ?, status = ?, payload_json = ?, error_message = ?, sent_at = ?, updated_at = ?
		WHERE id = ?`,
		notification.Channel, notification.Destination, string(notification.Status), string(payload), notification.ErrorMessage, emptyTime(notification.SentAt), db.NowString(notification.UpdatedAt), notification.ID,
	)
	return store.NotModified(result, err, "NOTIFICATION_NOT_FOUND", "通知记录不存在")
}

func (r *SQLRepository) FindByID(ctx context.Context, id string) (domain.Notification, error) {
	return scanNotification(r.store.QueryRow(ctx, notificationSelect+` WHERE id = ?`, id))
}

func (r *SQLRepository) ListByIncident(ctx context.Context, incidentID string, limit, offset int) ([]domain.Notification, int, error) {
	var total int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE incident_id = ?`, incidentID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count notifications: %w", err)
	}
	rows, err := r.store.Query(ctx, notificationSelect+` WHERE incident_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, incidentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var notifications []domain.Notification
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, total, rows.Err()
}

func (r *SQLRepository) List(ctx context.Context, channel string, status domain.Status, limit, offset int) ([]domain.Notification, int, error) {
	where := " WHERE 1=1"
	args := []any{}
	if channel != "" {
		where += " AND channel = ?"
		args = append(args, channel)
	}
	if status != "" {
		where += " AND status = ?"
		args = append(args, string(status))
	}
	var total int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM notifications`+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count notifications: %w", err)
	}
	query := notificationSelect + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := r.store.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var notifications []domain.Notification
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, total, rows.Err()
}

const notificationSelect = `
	SELECT id, incident_id, alert_id, channel, destination, status, payload_json, error_message, sent_at, created_at, updated_at
	FROM notifications`

type scanner interface {
	Scan(dest ...any) error
}

func scanNotification(row scanner) (domain.Notification, error) {
	var notification domain.Notification
	var status, payload, sentAt, createdAt, updatedAt string
	if err := row.Scan(&notification.ID, &notification.IncidentID, &notification.AlertID, &notification.Channel, &notification.Destination, &status, &payload, &notification.ErrorMessage, &sentAt, &createdAt, &updatedAt); err != nil {
		return domain.Notification{}, normalizeNotFound(err)
	}
	_ = json.Unmarshal([]byte(payload), &notification.Payload)
	notification.Status = domain.Status(status)
	notification.SentAt = parseOptionalTime(sentAt)
	notification.CreatedAt, _ = db.ParseTime(createdAt)
	notification.UpdatedAt, _ = db.ParseTime(updatedAt)
	return notification, nil
}

func emptyTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return db.NowString(value)
}

func parseOptionalTime(raw string) time.Time {
	value, _ := db.ParseTime(raw)
	return value
}

func normalizeNotFound(err error) error {
	if err == sql.ErrNoRows {
		return apperr.NotFound("NOTIFICATION_NOT_FOUND", "通知记录不存在")
	}
	return err
}
