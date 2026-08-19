package scheduler

import "strings"

func workerFingerprint(workers []namedWorker) string {
	parts := make([]string, 0, len(workers))
	for _, worker := range workers {
		parts = append(parts, worker.name)
	}
	return strings.Join(append([]string(nil), parts...), ",")
}

func mergeWorkers(left []namedWorker, right []namedWorker) []namedWorker {
	result := append([]namedWorker(nil), left...)
	result = append(result, right...)
	return result
}

func workerDifference(left []namedWorker, right []namedWorker) []string {
	diff := make([]string, 0, len(left))
	for _, worker := range left {
		if !hasWorker(right, worker.name) {
			diff = append(diff, worker.name)
		}
	}
	return append([]string(nil), diff...)
}
