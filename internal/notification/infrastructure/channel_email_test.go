package infrastructure

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/acme/signalforge/internal/notification/domain"
)

func TestChannelsHonorCanceledContext(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var failures []error
	for _, channel := range []interface {
		Send(context.Context, string, map[string]any) error
	}{NewEmailChannel(logger), NewLogChannel(logger)} {
		if err := channel.Send(ctx, "ops", nil); !errors.Is(err, context.Canceled) {
			failures = append(failures, err)
		}
	}
	if len(failures) != 0 {
		t.Fatalf("canceled channel errors = %v", failures)
	}
}

func TestDispatcherHonorsCanceledContext(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	dispatcher := NewDispatcher(NewEmailChannel(logger), NewLogChannel(logger), NewWebhookChannel())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := dispatcher.Send(ctx, domain.Notification{Channel: "email", Destination: "ops"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("dispatcher err=%v", err)
	}
}
