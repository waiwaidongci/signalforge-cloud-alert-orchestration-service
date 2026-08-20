package scheduler

import "database/sql"

func schedulerDependencyDB(r *Runner) *sql.DB {
	return r.deps.DB
}
