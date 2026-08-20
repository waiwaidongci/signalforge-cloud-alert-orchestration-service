package scheduler

import "database/sql"

func schedulerDependencyDB(r *Runner) *sql.DB {
	if r == nil || r.deps == nil {
		return nil
	}
	return r.deps.DB
}
