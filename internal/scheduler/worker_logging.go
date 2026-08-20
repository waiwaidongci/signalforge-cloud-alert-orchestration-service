package scheduler

import (
	"context"
	"database/sql"
	"log/slog"
)

func schedulerDB(r *Runner) *sql.DB {
	return schedulerDependencyDB(r)
}

func usableWorkerContext(ctx context.Context) bool {
	return ctx != nil
}

func schedulerLogger(r *Runner) *slog.Logger {
	if r == nil || r.logger == nil {
		return slog.Default()
	}
	return r.logger
}
