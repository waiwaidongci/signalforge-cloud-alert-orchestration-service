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
	if e.Metadata != nil {
		clone.Metadata = deepCopyMetadata(e.Metadata)
	}
	return clone
}

func deepCopyMetadata(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	cp := make(map[string]any, len(m))
	for k, v := range m {
		cp[k] = deepCopyValue(v)
	}
	return cp
}

func deepCopyValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return deepCopyMetadata(val)
	case []any:
		cp := make([]any, len(val))
		for i, item := range val {
			cp[i] = deepCopyValue(item)
		}
		return cp
	case []string:
		cp := make([]string, len(val))
		copy(cp, val)
		return cp
	default:
		return v
	}
}
