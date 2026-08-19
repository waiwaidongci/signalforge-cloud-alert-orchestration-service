package domain

import "time"

type EventType string

const (
	EventReceived     EventType = "received"
	EventDeduped      EventType = "deduped"
	EventAggregated   EventType = "aggregated"
	EventNotified     EventType = "notified"
	EventAcknowledged EventType = "acknowledged"
	EventResolved     EventType = "resolved"
	EventClosed       EventType = "closed"
	EventSuppressed   EventType = "suppressed"
)

type TimelineEvent struct {
	ID         string         `json:"id"`
	IncidentID string         `json:"incident_id"`
	AlertID    string         `json:"alert_id,omitempty"`
	EventType  EventType      `json:"event_type"`
	Actor      string         `json:"actor,omitempty"`
	Message    string         `json:"message"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	OccurredAt time.Time      `json:"occurred_at"`
	CreatedAt  time.Time      `json:"created_at"`
}
