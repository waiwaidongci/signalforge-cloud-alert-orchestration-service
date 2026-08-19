package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/acme/signalforge/internal/incident/domain"
	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/db"
	"github.com/acme/signalforge/internal/shared/store"
)

type TimelineSQLRepository struct {
	store *store.Store
}

func NewTimelineSQLRepository(raw *sql.DB) *TimelineSQLRepository {
	return &TimelineSQLRepository{store: store.New(raw)}
}

func (r *TimelineSQLRepository) Append(ctx context.Context, event domain.TimelineEvent) error {
	metadata, _ := json.Marshal(event.Metadata)
	_, err := r.store.Exec(ctx, `
		INSERT INTO timeline_events (id, incident_id, alert_id, event_type, actor, message, metadata_json, occurred_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.IncidentID, event.AlertID, string(event.EventType), event.Actor, event.Message, string(metadata), db.NowString(event.OccurredAt), db.NowString(event.CreatedAt),
	)
	if err != nil {
		return fmt.Errorf("append timeline event: %w", err)
	}
	return nil
}

func (r *TimelineSQLRepository) List(ctx context.Context, incidentID string, limit, offset int) ([]domain.TimelineEvent, int, error) {
	var total int
	if err := r.store.QueryRow(ctx, `SELECT COUNT(*) FROM timeline_events WHERE incident_id = ?`, incidentID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count timeline: %w", err)
	}
	rows, err := r.store.Query(ctx, timelineSelect+` WHERE incident_id = ? ORDER BY occurred_at DESC LIMIT ? OFFSET ?`, incidentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var events []domain.TimelineEvent
	for rows.Next() {
		event, err := scanTimeline(rows)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, event)
	}
	return events, total, rows.Err()
}

const timelineSelect = `
	SELECT id, incident_id, alert_id, event_type, actor, message, metadata_json, occurred_at, created_at
	FROM timeline_events`

func scanTimeline(row scanner) (domain.TimelineEvent, error) {
	var event domain.TimelineEvent
	var eventType, metadata, occurredAt, createdAt string
	if err := row.Scan(&event.ID, &event.IncidentID, &event.AlertID, &eventType, &event.Actor, &event.Message, &metadata, &occurredAt, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.TimelineEvent{}, apperr.NotFound("TIMELINE_NOT_FOUND", "时间线事件不存在")
		}
		return domain.TimelineEvent{}, err
	}
	_ = json.Unmarshal([]byte(metadata), &event.Metadata)
	event.EventType = domain.EventType(eventType)
	event.OccurredAt, _ = db.ParseTime(occurredAt)
	event.CreatedAt, _ = db.ParseTime(createdAt)
	return event, nil
}
