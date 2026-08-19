package scheduler

import (
	"context"
	"sync"
)

type namedWorker struct {
	name string
	run  func(context.Context)
}

func runWorkers(ctx context.Context, workers []namedWorker) map[string]string {
	collector := newResultCollector()
	var wg sync.WaitGroup
	for _, worker := range workers {
		wg.Add(1)
		go func(item namedWorker) {
			defer wg.Done()
			item.run(ctx)
			collector.Set(item.name)
		}(worker)
	}
	wg.Wait()
	return collector.Snapshot()
}
