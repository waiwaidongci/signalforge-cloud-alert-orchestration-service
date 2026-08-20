package infrastructure

import (
	"context"
	"log/slog"
)

type EmailChannel struct {
	logger *slog.Logger
}

func NewEmailChannel(logger *slog.Logger) *EmailChannel {
	return &EmailChannel{logger: logger}
}

func (c *EmailChannel) Name() string {
	return "email"
}

func (c *EmailChannel) Send(ctx context.Context, destination string, payload map[string]any) error {
	if err := channelContextError(ctx); err != nil {
		return err
	}
	if destination == "" {
		return context.Canceled
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	// This is a deterministic mail simulation kept explicit for local development.
	c.logger.Info("email notification simulated", "destination", destination, "payload", payload)
	return nil
}
