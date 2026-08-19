package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acme/signalforge/internal/app"
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
	handler := app.NewRouter(deps, app.RouterOptions{
		AuthToken:      cfg.Auth.PlaceholderToken,
		MaxBodyBytes:   cfg.Server.MaxBodyBytes,
		RequestTimeout: cfg.Server.RequestTimeout.Value(),
		RateLimit:      100,
	})

	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout.Value(),
		WriteTimeout: cfg.Server.WriteTimeout.Value(),
		IdleTimeout:  cfg.Server.IdleTimeout.Value(),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("signalforge server started", "addr", cfg.Server.Addr)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout.Value())
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	logger.Info("server stopped")
	return nil
}

func logStart(logger *slog.Logger, addr string, started time.Time) {
	logger.Info("listening", "addr", addr, "ready_in", time.Since(started).String())
}
