package infrastructure

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/acme/signalforge/internal/notification/domain"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestDispatcherUnknownChannelWrapsSentinel(t *testing.T) {
	dispatcher := NewDispatcher(NewLogChannel(discardLogger()))
	err := dispatcher.Send(context.Background(), domain.Notification{Channel: "missing"})
	if !errors.Is(err, domain.ErrUnknownChannel) {
		t.Fatalf("expected ErrUnknownChannel in chain, got %v", err)
	}
}

func TestDispatcherMissingDestinationWrapsSentinel(t *testing.T) {
	dispatcher := NewDispatcher(NewLogChannel(discardLogger()))
	err := dispatcher.Send(context.Background(), domain.Notification{Channel: "log"})
	if !errors.Is(err, domain.ErrDestinationRequired) {
		t.Fatalf("expected ErrDestinationRequired in chain, got %v", err)
	}
}
