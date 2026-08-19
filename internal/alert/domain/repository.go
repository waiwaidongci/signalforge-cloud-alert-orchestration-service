package domain

import (
	"context"
	"time"

	"github.com/acme/signalforge/internal/shared/severity"
)

type ListFilter struct {
	SourceID    string
	Fingerprint string
	Status      Status
	Severity    string
	IncidentID  string
	Query       string
	From        time.Time
	To          time.Time
}

type Repository interface {
	Create(ctx context.Context, alert Alert) error
	Update(ctx context.Context, alert Alert) error
	FindByID(ctx context.Context, id string) (Alert, error)
	FindBySourceExternalID(ctx context.Context, sourceID, externalID string) (Alert, error)
	FindByFingerprint(ctx context.Context, fingerprint string, limit int) ([]Alert, error)
	List(ctx context.Context, filter ListFilter, limit, offset int) ([]Alert, int, error)
	UpdateIncidentID(ctx context.Context, ids []string, incidentID string) error
	BatchUpdateStatus(ctx context.Context, ids []string, status Status, at time.Time) error
	CountByIncident(ctx context.Context, incidentID string, status Status) (int, error)
	HighestActiveSeverity(ctx context.Context) (severity.Severity, error)
}
