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
	return DispatchResult{ID: job.ID}
}
