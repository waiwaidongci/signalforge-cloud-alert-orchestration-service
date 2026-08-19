package infrastructure

import (
	"context"
	"log/slog"
)

type LogChannel struct {
	logger *slog.Logger
}

func NewLogChannel(logger *slog.Logger) *LogChannel {
	return &LogChannel{logger: logger}
}

func (c *LogChannel) Name() string {
	return "log"
}

func (c *LogChannel) Send(_ context.Context, destination string, payload map[string]any) error {
	c.logger.Info("notification log channel", "destination", destination, "payload", payload)
	return nil
}
