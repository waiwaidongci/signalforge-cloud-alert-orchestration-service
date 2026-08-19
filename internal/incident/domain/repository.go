package domain

import "context"
import "time"

type ListFilter struct {
	Fingerprint string
	Status      Status
	Severity    string
	Query       string
}

type Repository interface {
	Create(ctx context.Context, incident Incident) error
	Update(ctx context.Context, incident Incident) error
	FindByID(ctx context.Context, id string) (Incident, error)
	FindByFingerprint(ctx context.Context, fingerprint string) (Incident, error)
	List(ctx context.Context, filter ListFilter, limit, offset int) ([]Incident, int, error)
	IncrementAlertCount(ctx context.Context, id string, lastSeenAt time.Time) error
}

type TimelineRepository interface {
	Append(ctx context.Context, event TimelineEvent) error
	List(ctx context.Context, incidentID string, limit, offset int) ([]TimelineEvent, int, error)
}
