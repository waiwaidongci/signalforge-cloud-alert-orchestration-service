package logging

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/acme/signalforge/internal/shared/config"
)

// New builds a structured JSON or text slog.Logger.
func New(cfg config.LogConfig) (*slog.Logger, error) {
	level := slog.LevelInfo
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	options := &slog.HandlerOptions{Level: level}
	if cfg.Format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stdout, options)), nil
	}
	return slog.New(slog.NewTextHandler(os.Stdout, options)), nil
}

func WithRequestID(logger *slog.Logger, requestID string) *slog.Logger {
	return logger.With("request_id", requestID)
}

func PrettyJSON(w io.Writer, value any) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(value)
}
