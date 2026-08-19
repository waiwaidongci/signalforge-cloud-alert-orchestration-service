package domain

import (
	"time"

	"github.com/acme/signalforge/internal/shared/severity"
)

type Status string

const (
	StatusFiring     Status = "firing"
	StatusResolved   Status = "resolved"
	StatusSuppressed Status = "suppressed"
)

type Alert struct {
	ID             string            `json:"id"`
	SourceID       string            `json:"source_id"`
	ExternalID     string            `json:"external_id"`
	Fingerprint    string            `json:"fingerprint"`
	Resource       string            `json:"resource"`
	Severity       severity.Severity `json:"severity"`
	Status         Status            `json:"status"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Labels         map[string]string `json:"labels"`
	Annotations    map[string]string `json:"annotations"`
	ReceivedAt     time.Time         `json:"received_at"`
	LastOccurredAt time.Time         `json:"last_occurred_at"`
	DedupCount     int               `json:"dedup_count"`
	IncidentID     string            `json:"incident_id"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func (a Alert) Service() string {
	if value := a.Labels["service"]; value != "" {
		return value
	}
	return a.Resource
}

func (a Alert) Environment() string {
	return a.Labels["environment"]
}

func (a Alert) Target() (interface {
	Service() string
}, error) {
	return nil, nil
}
