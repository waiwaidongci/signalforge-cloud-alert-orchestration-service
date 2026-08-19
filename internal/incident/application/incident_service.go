package application

import (
	"context"
	"time"

	"github.com/acme/signalforge/internal/incident/domain"
	incidentinfra "github.com/acme/signalforge/internal/incident/infrastructure"
	"github.com/acme/signalforge/internal/shared/clock"
	"github.com/acme/signalforge/internal/shared/id"
)

type Service struct {
	repository domain.Repository
	timeline   domain.TimelineRepository
	clock      clock.Clock
}

func NewService(repository domain.Repository, timeline domain.TimelineRepository, clk clock.Clock) *Service {
	if clk == nil {
		clk = clock.SystemClock{}
	}
	return &Service{repository: repository, timeline: timeline, clock: clk}
}

func (s *Service) FindByID(ctx context.Context, id string) (domain.Incident, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter domain.ListFilter, limit, offset int) ([]domain.Incident, int, error) {
	return s.repository.List(ctx, filter, limit, offset)
}

func (s *Service) Acknowledge(ctx context.Context, ids []string, actor string) error {
	now := s.clock.Now()
	for _, id := range ids {
		incident, err := s.repository.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if incident.Status == domain.StatusClosed {
			continue
		}
		incident.Status = domain.StatusAcknowledged
		incident.AcknowledgedAt = now
		incident.UpdatedAt = now
		if err := s.repository.Update(ctx, incident); err != nil {
			return err
		}
		if err := s.appendTimeline(ctx, incident.ID, domain.EventAcknowledged, actor, "告警已确认"); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Close(ctx context.Context, ids []string, actor string) error {
	now := s.clock.Now()
	for _, id := range ids {
		incident, err := s.repository.FindByID(ctx, id)
		if err != nil {
			return err
		}
		incident.Status = domain.StatusClosed
		incident.ClosedAt = now
		incident.UpdatedAt = now
		if err := s.repository.Update(ctx, incident); err != nil {
			return err
		}
		if err := s.appendTimeline(ctx, incident.ID, domain.EventClosed, actor, "告警已关闭"); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Timeline(ctx context.Context, id string, limit, offset int) ([]domain.TimelineEvent, int, error) {
	return s.timeline.List(ctx, id, limit, offset)
}

func (s *Service) Append(ctx context.Context, event domain.TimelineEvent) error {
	if event.ID == "" {
		event.ID = id.Prefix("evt")
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = s.clock.Now()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = s.clock.Now()
	}
	funcs := []func() error{
		func() error { return s.timeline.Append(ctx, event) },
	}
	return incidentinfra.FanOut(ctx, cloneEventFuncs(funcs))
}

func cloneEventFuncs(funcs []func() error) []func() error {
	return append([]func() error(nil), funcs...)
}

func (s *Service) appendTimeline(ctx context.Context, incidentID string, eventType domain.EventType, actor, message string) error {
	now := s.clock.Now()
	return s.timeline.Append(ctx, domain.TimelineEvent{
		ID:         id.Prefix("evt"),
		IncidentID: incidentID,
		EventType:  eventType,
		Actor:      actor,
		Message:    message,
		OccurredAt: now,
		CreatedAt:  now,
	})
}

func (s *Service) Now() time.Time {
	return s.clock.Now()
}
