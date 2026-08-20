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

func (e TimelineEvent) Clone() TimelineEvent {
	clone := e
	clone.Metadata = cloneTimelineMetadata(e.Metadata)
	return clone
}

func cloneTimelineMetadata(metadata map[string]any) map[string]any {
	if metadata == nil {
		return nil
	}
	clone := make(map[string]any, len(metadata))
	for key, value := range metadata {
		clone[key] = cloneTimelineValue(value)
	}
	return clone
}

func cloneTimelineValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return cloneTimelineMetadata(typed)
	case []string:
		return append([]string(nil), typed...)
	default:
		return value
	}
}
