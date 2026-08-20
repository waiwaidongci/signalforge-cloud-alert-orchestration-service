package scheduler

import (
	"context"

	incidentdomain "github.com/acme/signalforge/internal/incident/domain"
)

func (r *Runner) runRecovery(ctx context.Context) {
	db := schedulerDB(r)
	if db == nil || !usableWorkerContext(ctx) {
		return
	}
	incidents, _, err := r.deps.Incidents.List(ctx, incidentdomain.ListFilter{Status: incidentdomain.StatusActive}, 100, 0)
	if err != nil {
		schedulerLogger(r).Warn("list active incidents failed", "error", err)
		return
	}
	for _, incident := range incidents {
		var firing int
		row := r.deps.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE incident_id = ? AND status = ?`, incident.ID, "firing")
		if err := row.Scan(&firing); err != nil {
			schedulerLogger(r).Warn("count firing alerts failed", "incident_id", incident.ID, "error", err)
			continue
		}
		if firing > 0 {
			continue
		}
		if err := r.deps.Incidents.Close(ctx, []string{incident.ID}, "scheduler"); err != nil {
			schedulerLogger(r).Warn("auto close incident failed", "incident_id", incident.ID, "error", err)
		}
	}
}
