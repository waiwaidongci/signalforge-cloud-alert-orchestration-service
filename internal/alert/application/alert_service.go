package application

import (
	"context"
	"log/slog"

	alertdomain "github.com/acme/signalforge/internal/alert/domain"
	incidentdomain "github.com/acme/signalforge/internal/incident/domain"
	"github.com/acme/signalforge/internal/shared/clock"
	idgen "github.com/acme/signalforge/internal/shared/id"
	"github.com/acme/signalforge/internal/shared/metrics"
)

type Service struct {
	alerts    AlertRepository
	incidents IncidentRepository
	timeline  TimelineWriter
	clock     clock.Clock
	logger    *slog.Logger
	metrics   *metrics.Metrics
}

func NewService(alerts AlertRepository, incidents IncidentRepository, timeline TimelineWriter, clk clock.Clock, logger *slog.Logger, metrics *metrics.Metrics) *Service {
	if clk == nil {
		clk = clock.SystemClock{}
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{alerts: alerts, incidents: incidents, timeline: timeline, clock: clk, logger: logger, metrics: metrics}
}

func (s *Service) List(ctx context.Context, filter alertdomain.ListFilter, limit, offset int) ([]alertdomain.Alert, int, error) {
	return s.alerts.List(ctx, filter, limit, offset)
}

func (s *Service) Get(ctx context.Context, id string) (alertdomain.Alert, error) {
	return s.alerts.FindByID(ctx, id)
}

func (s *Service) Acknowledge(ctx context.Context, ids []string, actor string) error {
	now := s.clock.Now()
	if err := s.alerts.BatchUpdateStatus(ctx, ids, alertdomain.StatusAcknowledged, now); err != nil {
		return err
	}
	for _, id := range ids {
		alert, err := s.alerts.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if alert.IncidentID != "" {
			if err := s.timeline.Append(ctx, incidentdomain.TimelineEvent{
				ID:         idgen.Prefix("evt"),
				IncidentID: alert.IncidentID,
				AlertID:    alert.ID,
				EventType:  incidentdomain.EventAcknowledged,
				Actor:      actor,
				Message:    "告警已批量确认",
				OccurredAt: now,
				CreatedAt:  now,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) Close(ctx context.Context, ids []string, actor string) error {
	now := s.clock.Now()
	if err := s.alerts.BatchUpdateStatus(ctx, ids, alertdomain.StatusResolved, now); err != nil {
		return err
	}
	for _, id := range ids {
		alert, err := s.alerts.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if alert.IncidentID != "" {
			if err := s.timeline.Append(ctx, incidentdomain.TimelineEvent{
				ID:         idgen.Prefix("evt"),
				IncidentID: alert.IncidentID,
				AlertID:    alert.ID,
				EventType:  incidentdomain.EventClosed,
				Actor:      actor,
				Message:    "告警已批量关闭",
				OccurredAt: now,
				CreatedAt:  now,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
