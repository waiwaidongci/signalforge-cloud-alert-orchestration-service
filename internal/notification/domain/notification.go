package domain

import "time"

type Status string

const (
	StatusPending Status = "pending"
	StatusSent    Status = "sent"
	StatusFailed  Status = "failed"
)

type Notification struct {
	ID           string         `json:"id"`
	IncidentID   string         `json:"incident_id"`
	AlertID      string         `json:"alert_id"`
	Channel      string         `json:"channel"`
	Destination  string         `json:"destination"`
	Status       Status         `json:"status"`
	Payload      map[string]any `json:"payload"`
	ErrorMessage string         `json:"error_message,omitempty"`
	SentAt       time.Time      `json:"sent_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
