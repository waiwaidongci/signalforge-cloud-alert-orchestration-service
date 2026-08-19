package scheduler

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/acme/signalforge/internal/app"
	"github.com/acme/signalforge/internal/shared/metrics"
)

type Runner struct {
	deps     *app.Dependencies
	logger   *slog.Logger
	metrics  *metrics.Metrics
	interval time.Duration
	jitter   time.Duration
}

func NewRunner(deps *app.Dependencies, logger *slog.Logger, metrics *metrics.Metrics, interval, jitter time.Duration) *Runner {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	if jitter <= 0 {
		jitter = time.Second
	}
	return &Runner{deps: deps, logger: logger, metrics: metrics, interval: interval, jitter: jitter}
}

func (r *Runner) Start(ctx context.Context) {
	go func() {
		for {
			delay := r.interval
			if r.jitter > 0 {
				delay += time.Duration(rand.Int63n(int64(r.jitter)))
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				r.runOnce(ctx)
			}
		}
	}()
}

func (r *Runner) runOnce(ctx context.Context) {
	if r.metrics != nil {
		r.metrics.SchedulerRuns.Inc()
	}
	r.logger.Info("scheduler loop started")
	start := time.Now()
	workers := []struct {
		name string
		run  func(context.Context)
	}{
		{"recovery", r.runRecovery},
		{"escalation", r.runEscalation},
		{"silence_expiry", r.runSilenceExpiry},
	}
	for _, worker := range workers {
		worker.run(ctx)
	}
	r.logger.Info("scheduler loop finished", "duration", time.Since(start).String())
}
