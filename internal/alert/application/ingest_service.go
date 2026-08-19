package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	alertdomain "github.com/acme/signalforge/internal/alert/domain"
	dedupdomain "github.com/acme/signalforge/internal/dedup/domain"
	incidentdomain "github.com/acme/signalforge/internal/incident/domain"
	notificationdomain "github.com/acme/signalforge/internal/notification/domain"
	"github.com/acme/signalforge/internal/shared/clock"
	"github.com/acme/signalforge/internal/shared/id"
	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/shared/metrics"
	"github.com/acme/signalforge/internal/shared/severity"
)

type IngestService struct {
	alerts       AlertRepository
	sources      SourceRepository
	fingerprints FingerprintResolver
	silences     SilencePolicy
	routes       RouteResolver
	notifier     Notifier
	incidents    IncidentRepository
	timeline     TimelineWriter
	clock        clock.Clock
	logger       *slog.Logger
	metrics      *metrics.Metrics
}

func NewIngestService(
	alerts AlertRepository,
	sources SourceRepository,
	fingerprints FingerprintResolver,
	silences SilencePolicy,
	routes RouteResolver,
	notifier Notifier,
	incidents IncidentRepository,
	timeline TimelineWriter,
	clk clock.Clock,
	logger *slog.Logger,
	metrics *metrics.Metrics,
) *IngestService {
	if clk == nil {
		clk = clock.SystemClock{}
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &IngestService{
		alerts: alerts, sources: sources, fingerprints: fingerprints, silences: silences,
		routes: routes, notifier: notifier, incidents: incidents, timeline: timeline,
		clock: clk, logger: logger, metrics: metrics,
	}
}

func (s *IngestService) Ingest(ctx context.Context, input alertdomain.IngestInput) (alertdomain.IngestResult, error) {
	now := s.clock.Now()
	if input.ReceivedAt.IsZero() {
		input.ReceivedAt = now
	}
	source, err := s.sources.FindByID(ctx, input.SourceID)
	if err != nil {
		return alertdomain.IngestResult{}, err
	}
	if !source.Enabled {
		return alertdomain.IngestResult{}, fmt.Errorf("source %s is disabled", source.ID)
	}

	fingerprint, err := s.fingerprints.Resolve(ctx, dedupdomain.Input{
		SourceID: input.SourceID,
		Resource: input.Resource,
		Severity: severity.Severity(input.Severity),
		Title:    input.Title,
		Labels:   input.Labels,
	})
	if err != nil {
		return alertdomain.IngestResult{}, err
	}

	if existing, err := s.alerts.FindBySourceExternalID(ctx, input.SourceID, input.ExternalID); err == nil {
		return s.deduplicateExisting(ctx, existing, input, now)
	}

	recent, _ := s.alerts.FindByFingerprint(ctx, fingerprint, 5)
	if len(recent) > 0 {
		primary := recent[0]
		primary.DedupCount++
		primary.LastOccurredAt = now
		primary.Status = alertdomain.StatusFiring
		primary.UpdatedAt = now
		if err := s.alerts.Update(ctx, primary); err != nil {
			return alertdomain.IngestResult{}, err
		}
		if primary.IncidentID == "" {
			incident, _ := s.ensureIncident(ctx, primary)
			primary.IncidentID = incident.ID
			_ = s.alerts.Update(ctx, primary)
		} else {
			_ = s.incidents.IncrementAlertCount(ctx, primary.IncidentID, now)
		}
		s.appendTimeline(ctx, primary.IncidentID, primary.ID, incidentdomain.EventDeduped, "相同指纹告警已合并")
		if s.metrics != nil {
			s.metrics.Alerts.WithLabelValues(primary.Severity.String(), string(primary.Status)).Inc()
		}
		return alertdomain.IngestResult{
			AlertID: primary.ID, IncidentID: primary.IncidentID, Created: false, Deduped: true, Suppressed: false,
			AlertStatus: primary.Status, IncidentStatus: incidentStatus(ctx, s.incidents, primary.IncidentID),
		}, nil
	}

	alert := alertdomain.Alert{
		ID:             id.Prefix("alr"),
		SourceID:       input.SourceID,
		ExternalID:     input.ExternalID,
		Fingerprint:    fingerprint,
		Resource:       input.Resource,
		Severity:       severity.Severity(input.Severity),
		Status:         alertdomain.StatusFiring,
		Title:          input.Title,
		Description:    input.Description,
		Labels:         input.Labels,
		Annotations:    input.Annotations,
		ReceivedAt:     input.ReceivedAt,
		LastOccurredAt: input.ReceivedAt,
		DedupCount:     1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	target := matcher.Target{SourceID: source.ID, Service: alert.Service(), Env: alert.Environment(), Labels: alert.Labels, Severity: alert.Severity}
	silenced := s.silencedBy(ctx, target)
	if silenced != "" {
		alert.Status = alertdomain.StatusSuppressed
	}
	if err := s.alerts.Create(ctx, alert); err != nil {
		return alertdomain.IngestResult{}, err
	}
	incident, err := s.ensureIncident(ctx, alert)
	if err != nil {
		return alertdomain.IngestResult{}, err
	}
	alert.IncidentID = incident.ID
	_ = s.alerts.Update(ctx, alert)
	s.appendTimeline(ctx, incident.ID, alert.ID, incidentdomain.EventReceived, "告警已接收")
	if alert.Status != alertdomain.StatusSuppressed {
		s.routeAndNotify(ctx, alert, incident)
	} else {
		s.appendTimeline(ctx, incident.ID, alert.ID, incidentdomain.EventSuppressed, "告警已被静默规则抑制")
	}
	if s.metrics != nil {
		s.metrics.Alerts.WithLabelValues(alert.Severity.String(), string(alert.Status)).Inc()
	}
	return alertdomain.IngestResult{
		AlertID: alert.ID, IncidentID: incident.ID, Created: true, Deduped: false, Suppressed: alert.Status == alertdomain.StatusSuppressed,
		AlertStatus: alert.Status, IncidentStatus: string(incident.Status),
	}, nil
}

func (s *IngestService) deduplicateExisting(ctx context.Context, existing alertdomain.Alert, input alertdomain.IngestInput, now time.Time) (alertdomain.IngestResult, error) {
	existing.DedupCount++
	existing.LastOccurredAt = input.ReceivedAt
	existing.UpdatedAt = now
	if existing.Status == alertdomain.StatusResolved {
		existing.Status = alertdomain.StatusFiring
	}
	if err := s.alerts.Update(ctx, existing); err != nil {
		return alertdomain.IngestResult{}, err
	}
	if existing.IncidentID == "" {
		incident, _ := s.ensureIncident(ctx, existing)
		existing.IncidentID = incident.ID
		_ = s.alerts.Update(ctx, existing)
	}
	s.appendTimeline(ctx, existing.IncidentID, existing.ID, incidentdomain.EventDeduped, "相同 external_id 告警已幂等合并")
	return alertdomain.IngestResult{AlertID: existing.ID, IncidentID: existing.IncidentID, Created: false, Deduped: true, Suppressed: false, AlertStatus: existing.Status, IncidentStatus: incidentStatus(ctx, s.incidents, existing.IncidentID)}, nil
}

func (s *IngestService) ensureIncident(ctx context.Context, alert alertdomain.Alert) (incidentdomain.Incident, error) {
	if alert.IncidentID != "" {
		if incident, err := s.incidents.FindByFingerprint(ctx, alert.Fingerprint); err == nil {
			return incident, nil
		}
	}
	existing, err := s.incidents.FindByFingerprint(ctx, alert.Fingerprint)
	if err == nil {
		existing.AlertCount++
		existing.LastSeenAt = alert.LastOccurredAt
		existing.UpdatedAt = s.clock.Now()
		if err := s.incidents.Update(ctx, existing); err != nil {
			return incidentdomain.Incident{}, err
		}
		return existing, nil
	}
	now := s.clock.Now()
	incident := incidentdomain.Incident{
		ID:          id.Prefix("inc"),
		Fingerprint: alert.Fingerprint,
		Title:       alert.Title,
		Severity:    alert.Severity,
		Status:      incidentdomain.StatusActive,
		SourceID:    alert.SourceID,
		AlertCount:  1,
		FirstSeenAt: alert.ReceivedAt,
		LastSeenAt:  alert.ReceivedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.incidents.Create(ctx, incident); err != nil {
		return incidentdomain.Incident{}, err
	}
	return incident, nil
}

func (s *IngestService) silencedBy(ctx context.Context, target matcher.Target) string {
	active, err := s.silences.ActiveSilences(ctx)
	if err != nil {
		s.logger.Warn("list active silences failed", "error", err)
		return ""
	}
	for i := range active {
		if active[i].Active(s.clock.Now()) && active[i].Matcher.Matches(target) {
			return active[i].ID
		}
	}
	highest, _ := s.alerts.HighestActiveSeverity(ctx)
	if suppressed, err := s.silences.IsSuppressed(ctx, target, target.Severity, severity.HigherOrEqual(highest, target.Severity)); err == nil && suppressed {
		return "suppression"
	}
	return ""
}

func (s *IngestService) routeAndNotify(ctx context.Context, alert alertdomain.Alert, incident incidentdomain.Incident) {
	channels, err := s.routes.Resolve(ctx, matcher.Target{SourceID: alert.SourceID, Service: alert.Service(), Env: alert.Environment(), Labels: alert.Labels, Severity: alert.Severity})
	if err != nil {
		s.logger.Warn("resolve routes failed", "error", err)
		return
	}
	if len(channels) == 0 {
		return
	}
	for _, channel := range channels {
		notification := notificationdomain.Notification{
			IncidentID: incident.ID, AlertID: alert.ID, Channel: channel.Channel, Destination: channel.Destination,
			Payload: map[string]any{"title": alert.Title, "severity": alert.Severity.String(), "resource": alert.Resource},
		}
		if err := s.notifier.Notify(ctx, notification); err != nil {
			s.logger.Warn("notify failed", "channel", channel.Channel, "error", err)
			continue
		}
		s.appendTimeline(ctx, incident.ID, alert.ID, incidentdomain.EventNotified, "通知已发送到 "+channel.Channel)
	}
}

func (s *IngestService) appendTimeline(ctx context.Context, incidentID, alertID string, eventType incidentdomain.EventType, message string) {
	now := s.clock.Now()
	if err := s.timeline.Append(ctx, incidentdomain.TimelineEvent{
		ID: id.Prefix("evt"), IncidentID: incidentID, AlertID: alertID, EventType: eventType, Message: message, OccurredAt: now, CreatedAt: now,
	}); err != nil {
		s.logger.Warn("append timeline failed", "incident_id", incidentID, "error", err)
	}
}

func incidentStatus(ctx context.Context, repo IncidentRepository, id string) string {
	if id == "" {
		return ""
	}
	incident, err := repo.FindByID(ctx, id)
	if err != nil {
		return ""
	}
	return string(incident.Status)
}
