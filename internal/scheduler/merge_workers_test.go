package scheduler

import (
	"sync"
	"testing"
)

func TestMergeWorkersDoesNotMutateLeft(t *testing.T) {
	left := make([]namedWorker, 1, 2)
	left[0] = namedWorker{name: "one"}
	right := []namedWorker{{name: "two"}}
	before := append([]namedWorker(nil), left...)
	got := mergeWorkers(left, right)
	if len(got) != 2 {
		t.Fatalf("expected 2 workers, got %d", len(got))
	}
	if len(left) != len(before) || left[0].name != before[0].name {
		t.Fatalf("left was mutated: before=%+v after=%+v", before, left)
	}
	if left[:2][1].name != "" {
		t.Fatalf("left backing array was mutated: %+v", left[:2])
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_ = mergeWorkers(left, right)
		}()
	}
	close(start)
	wg.Wait()
}
