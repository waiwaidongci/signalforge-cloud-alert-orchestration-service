package domain

import (
	"context"
	"time"
)

type UseCase interface {
	Ingest(ctx context.Context, input IngestInput) (IngestResult, error)
	List(ctx context.Context, filter ListFilter, limit, offset int) ([]Alert, int, error)
	Acknowledge(ctx context.Context, ids []string, actor string) error
	Close(ctx context.Context, ids []string, actor string) error
}

type IngestInput struct {
	SourceID    string
	ExternalID  string
	Resource    string
	Severity    string
	Title       string
	Description string
	Labels      map[string]string
	Annotations map[string]string
	ReceivedAt  time.Time
}

type IngestResult struct {
	AlertID        string
	IncidentID     string
	Created        bool
	Deduped        bool
	Suppressed     bool
	AlertStatus    Status
	IncidentStatus string
}
