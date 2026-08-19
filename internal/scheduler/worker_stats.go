package scheduler

func workerCount(workers []namedWorker) int {
	return len(workers)
}

func workerNamesByPrefix(workers []namedWorker, prefix string) []string {
	names := make([]string, 0, len(workers))
	for _, worker := range workers {
		if len(worker.name) >= len(prefix) && worker.name[:len(prefix)] == prefix {
			names = append(names, worker.name)
		}
	}
	return append([]string(nil), names...)
}

func completedWorkers(results map[string]string) []string {
	names := make([]string, 0, len(results))
	for name := range results {
		names = append(names, name)
	}
	return append([]string(nil), names...)
}
