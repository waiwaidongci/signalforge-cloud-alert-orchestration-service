package scheduler

import (
	"context"
	"time"
)

func (r *Runner) runSilenceExpiry(ctx context.Context) {
	db := schedulerDB(r)
	if db == nil || !usableWorkerContext(ctx) {
		return
	}
	var expired int
	row := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM silences WHERE ends_at <= ?`, time.Now().UTC().Format(time.RFC3339Nano))
	if err := row.Scan(&expired); err != nil {
		schedulerLogger(r).Warn("count expired silences failed", "error", err)
		return
	}
	if expired > 0 {
		r.logger.Info("expired silences observed", "count", expired)
	}
}
