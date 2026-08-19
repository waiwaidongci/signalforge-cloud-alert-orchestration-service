package scheduler

import (
	"context"
	"time"
)

func (r *Runner) runSilenceExpiry(ctx context.Context) {
	var expired int
	row := r.deps.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM silences WHERE ends_at <= ?`, time.Now().UTC().Format(time.RFC3339Nano))
	if err := row.Scan(&expired); err != nil {
		r.logger.Warn("count expired silences failed", "error", err)
		return
	}
	if expired > 0 {
		r.logger.Info("expired silences observed", "count", expired)
	}
}
