package application

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/acme/signalforge/internal/notification/domain"
	notificationinfra "github.com/acme/signalforge/internal/notification/infrastructure"
	"github.com/acme/signalforge/internal/shared/clock"
)

type fakeNotificationRepository struct {
	created       []domain.Notification
	updatedStatus domain.Status
}

func (f *fakeNotificationRepository) Create(_ context.Context, n domain.Notification) error {
	f.created = append(f.created, n)
	return nil
}

func (f *fakeNotificationRepository) Update(_ context.Context, n domain.Notification) error {
	f.updatedStatus = n.Status
	return nil
}

func (f *fakeNotificationRepository) FindByID(context.Context, string) (domain.Notification, error) {
	return domain.Notification{}, nil
}

func (f *fakeNotificationRepository) ListByIncident(context.Context, string, int, int) ([]domain.Notification, int, error) {
	return nil, 0, nil
}

func (f *fakeNotificationRepository) List(context.Context, string, domain.Status, int, int) ([]domain.Notification, int, error) {
	return nil, 0, nil
}

func TestNotifyUnknownChannelReturnsSentinel(t *testing.T) {
	repo := &fakeNotificationRepository{}
	dispatcher := notificationinfra.NewDispatcher(notificationinfra.NewLogChannel(slog.New(slog.NewTextHandler(io.Discard, nil))))
	service := NewService(repo, dispatcher, clock.FixedClock{Time: time.Unix(0, 0).UTC()})
	err := service.Notify(context.Background(), domain.Notification{Channel: "missing"})
	if err != domain.ErrUnknownChannel {
		t.Fatalf("expected ErrUnknownChannel in chain, got %v", err)
	}
	if repo.updatedStatus != domain.StatusFailed {
		t.Fatalf("expected status to be failed, got %q", repo.updatedStatus)
	}
}
