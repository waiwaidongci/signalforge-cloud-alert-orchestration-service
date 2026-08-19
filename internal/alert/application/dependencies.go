package application

import (
	"context"
	"time"

	alertdomain "github.com/acme/signalforge/internal/alert/domain"
	dedupdomain "github.com/acme/signalforge/internal/dedup/domain"
	incidentdomain "github.com/acme/signalforge/internal/incident/domain"
	notificationdomain "github.com/acme/signalforge/internal/notification/domain"
	routingdomain "github.com/acme/signalforge/internal/routing/domain"
	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/shared/severity"
	silencedomain "github.com/acme/signalforge/internal/silence/domain"
	sourcedomain "github.com/acme/signalforge/internal/source/domain"
)

type AlertRepository interface {
	alertdomain.Repository
}

type SourceRepository interface {
	FindByID(ctx context.Context, id string) (sourcedomain.Source, error)
	FindByAPIKey(ctx context.Context, apiKey string) (sourcedomain.Source, error)
}

type FingerprintResolver interface {
	Resolve(ctx context.Context, input dedupdomain.Input) (string, error)
}

type SilencePolicy interface {
	ActiveSilences(ctx context.Context) ([]silencedomain.Silence, error)
	IsSuppressed(ctx context.Context, target matcher.Target, currentSeverity severity.Severity, highActive bool) (bool, error)
}

type RouteResolver interface {
	Resolve(ctx context.Context, target matcher.Target) ([]routingdomain.Channel, error)
}

type Notifier interface {
	Notify(ctx context.Context, notification notificationdomain.Notification) error
}

type IncidentRepository interface {
	FindByID(ctx context.Context, id string) (incidentdomain.Incident, error)
	FindByFingerprint(ctx context.Context, fingerprint string) (incidentdomain.Incident, error)
	Create(ctx context.Context, incident incidentdomain.Incident) error
	Update(ctx context.Context, incident incidentdomain.Incident) error
	IncrementAlertCount(ctx context.Context, id string, lastSeenAt time.Time) error
}

type TimelineWriter interface {
	Append(ctx context.Context, event incidentdomain.TimelineEvent) error
}
