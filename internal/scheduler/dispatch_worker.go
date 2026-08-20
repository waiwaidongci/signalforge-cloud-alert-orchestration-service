package scheduler

type DispatchJob struct {
	ID   string
	Fail bool
}
type DispatchResult struct {
	ID  string
	Err error
}

func dispatchFailure(job DispatchJob) error {
	if job.Fail {
		return ErrDispatchFailed
	}
	return nil
}

func dispatchResult(job DispatchJob) DispatchResult {
	return DispatchResult{ID: job.ID, Err: dispatchFailure(job)}
}

func executeDispatch(job DispatchJob) DispatchResult {
	return dispatchResult(job)
}
