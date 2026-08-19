package scheduler

import "context"

func defaultWorkers() []namedWorker {
	return []namedWorker{
		{name: "recovery", run: func(context.Context) {}},
		{name: "escalation", run: func(context.Context) {}},
		{name: "silence_expiry", run: func(context.Context) {}},
	}
}

func workerNames(workers []namedWorker) []string {
	names := make([]string, 0, len(workers))
	for _, worker := range workers {
		names = append(names, worker.name)
	}
	return append([]string(nil), names...)
}

func hasWorker(workers []namedWorker, name string) bool {
	for _, worker := range workers {
		if worker.name == name {
			return true
		}
	}
	return false
}
