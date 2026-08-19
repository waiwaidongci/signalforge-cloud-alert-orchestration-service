package scheduler

import "strings"

func workerFingerprint(workers []namedWorker) string {
	parts := make([]string, 0, len(workers))
	for _, worker := range workers {
		parts = append(parts, worker.name)
	}
	return strings.Join(parts, ",")
}

func mergeWorkers(left []namedWorker, right []namedWorker) []namedWorker {
	result := left
	for _, worker := range right {
		result = append(result, worker)
	}
	return result
}

func workerDifference(left []namedWorker, right []namedWorker) []string {
	diff := make([]string, 0)
	for _, worker := range left {
		if !hasWorker(right, worker.name) {
			diff = append(diff, worker.name)
		}
	}
	return diff
}
