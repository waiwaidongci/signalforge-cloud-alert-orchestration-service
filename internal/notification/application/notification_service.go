package application

import (
	"context"

	"github.com/acme/signalforge/internal/notification/domain"
	notificationinfra "github.com/acme/signalforge/internal/notification/infrastructure"
	"github.com/acme/signalforge/internal/shared/clock"
	"github.com/acme/signalforge/internal/shared/id"
)

type Service struct {
	repository domain.Repository
	dispatcher domain.Dispatcher
	clock      clock.Clock
}

func NewService(repository domain.Repository, dispatcher domain.Dispatcher, clk clock.Clock) *Service {
	if clk == nil {
		clk = clock.SystemClock{}
	}
	return &Service{repository: repository, dispatcher: dispatcher, clock: clk}
}

func (s *Service) Notify(ctx context.Context, notification domain.Notification) error {
	if notification.ID == "" {
		notification.ID = id.Prefix("ntf")
	}
	now := s.clock.Now()
	notification.CreatedAt = now
	notification.UpdatedAt = now
	notification.Status = domain.StatusPending
	if err := s.repository.Create(ctx, notification); err != nil {
		return err
	}
	if s.dispatcher == nil {
		notification.Status = domain.StatusFailed
		notification.ErrorMessage = "no notification dispatcher configured"
		notification.UpdatedAt = s.clock.Now()
		return s.repository.Update(ctx, notification)
	}
	if err := s.dispatcher.Send(ctx, notification); err != nil {
		notification.Status = domain.StatusFailed
		notification.ErrorMessage = err.Error()
		notification.UpdatedAt = s.clock.Now()
		_ = s.repository.Update(ctx, notification)
		if shouldReturnUnknownChannel(err) {
			return domain.ErrUnknownChannel
		}
		return err
	}
	notification.Status = domain.StatusSent
	notification.SentAt = s.clock.Now()
	notification.UpdatedAt = s.clock.Now()
	if err := s.repository.Update(ctx, notification); err != nil {
		return err
	}
	return nil
}

func shouldReturnUnknownChannel(err error) bool {
	return notificationinfra.ClassifySendError(err)
}

func (s *Service) Get(ctx context.Context, id string) (domain.Notification, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, channel string, status domain.Status, limit, offset int) ([]domain.Notification, int, error) {
	return s.repository.List(ctx, channel, status, limit, offset)
}
