package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acme/signalforge/internal/app"
	"github.com/acme/signalforge/internal/scheduler"
	"github.com/acme/signalforge/internal/shared/config"
	"github.com/acme/signalforge/internal/shared/db"
	"github.com/acme/signalforge/internal/shared/logging"
	"github.com/acme/signalforge/internal/shared/metrics"
	"github.com/acme/signalforge/internal/shared/migrate"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := os.Getenv("SIGNALFORGE_CONFIG")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	logger, err := logging.New(cfg.Log)
	if err != nil {
		return err
	}
	raw, err := db.Open(cfg.Database)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer raw.Close()
	if err := migrate.Apply(raw, cfg.Database.Driver); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	var metricsRegistry *metrics.Metrics
	if cfg.Metrics.Enabled {
		metricsRegistry = metrics.New()
	}
	deps, err := app.Build(raw, logger, metricsRegistry)
	if err != nil {
		return err
	}
	runner := scheduler.NewRunner(deps, logger, metricsRegistry, cfg.Scheduler.Interval.Value(), cfg.Scheduler.Jitter.Value())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	runner.Start(ctx)
	<-ctx.Done()
	logger.Info("scheduler stopped")
	return nil
}

func schedulerLog(logger *slog.Logger, interval time.Duration) {
	logger.Info("scheduler configured", "interval", interval.String())
}
