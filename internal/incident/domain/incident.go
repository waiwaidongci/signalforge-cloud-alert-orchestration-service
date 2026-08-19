package domain

import (
	"time"

	"github.com/acme/signalforge/internal/shared/severity"
)

type Status string

const (
	StatusActive       Status = "active"
	StatusAcknowledged Status = "acknowledged"
	StatusResolved     Status = "resolved"
	StatusClosed       Status = "closed"
)

type Incident struct {
	ID             string            `json:"id"`
	Fingerprint    string            `json:"fingerprint"`
	Title          string            `json:"title"`
	Severity       severity.Severity `json:"severity"`
	Status         Status            `json:"status"`
	SourceID       string            `json:"source_id"`
	AlertCount     int               `json:"alert_count"`
	FirstSeenAt    time.Time         `json:"first_seen_at"`
	LastSeenAt     time.Time         `json:"last_seen_at"`
	AcknowledgedAt time.Time         `json:"acknowledged_at,omitempty"`
	ClosedAt       time.Time         `json:"closed_at,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusActive, StatusAcknowledged, StatusResolved, StatusClosed:
		return true
	default:
		return false
	}
}
