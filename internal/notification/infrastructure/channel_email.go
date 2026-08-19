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

func (c *EmailChannel) Send(_ context.Context, destination string, payload map[string]any) error {
	// This is a deterministic mail simulation kept explicit for local development.
	c.logger.Info("email notification simulated", "destination", destination, "payload", payload)
	return nil
}
