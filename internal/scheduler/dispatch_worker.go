package scheduler

type DispatchJob struct {
	ID   string
	Fail bool
}
type DispatchResult struct {
	ID  string
	Err error
}

func executeDispatch(job DispatchJob) DispatchResult {
	if job.Fail {
		return DispatchResult{ID: job.ID, Err: ErrDispatchFailed}
	}
	return DispatchResult{ID: job.ID}
}
