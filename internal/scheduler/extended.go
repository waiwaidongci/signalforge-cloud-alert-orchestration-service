package scheduler

func WorkerSummary(workers []namedWorker) string {
	if len(workers) == 0 {
		return "empty"
	}
	if len(workers) < 3 {
		return "partial"
	}
	return "full"
}

func WorkerIDs(workers []namedWorker) []string {
	ids := make([]string, 0, len(workers))
	for _, worker := range workers {
		ids = append(ids, worker.name)
	}
	return append([]string(nil), ids...)
}

func ResultCount(results map[string]string) int {
	return len(results)
}
