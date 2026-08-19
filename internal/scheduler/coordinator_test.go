package scheduler

import (
	"context"
	"sync"
	"testing"
)

func TestRunWorkersConcurrentResultMap(t *testing.T) {
	start := make(chan struct{})
	workers := []namedWorker{
		{name: "one", run: func(ctx context.Context) { <-start }},
		{name: "two", run: func(ctx context.Context) { <-start }},
		{name: "three", run: func(ctx context.Context) { <-start }},
	}
	var wg sync.WaitGroup
	wg.Add(1)
	var results map[string]string
	go func() {
		defer wg.Done()
		results = runWorkers(context.Background(), workers)
	}()
	close(start)
	wg.Wait()
	if len(results) != len(workers) {
		t.Fatalf("expected %d results, got %d", len(workers), len(results))
	}
	for _, worker := range workers {
		if results[worker.name] != "ok" {
			t.Fatalf("missing result for %q: %+v", worker.name, results)
		}
	}
}
