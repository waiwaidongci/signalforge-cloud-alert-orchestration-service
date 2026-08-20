package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestDispatchBatchReturnsEveryResult(t *testing.T) {
	done := make(chan []DispatchResult, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			done <- DispatchBatch([]DispatchJob{{ID: "a"}, {ID: "b", Fail: true}})
		}()
	}
	close(start)
	go func() { wg.Wait(); close(done) }()
	select {
	case results := <-done:
		if !dispatchResultsComplete(results, 2) {
			t.Fatalf("want two results, got %#v", results)
		}
		failed := false
		for _, result := range results {
			failed = failed || result.Err != nil
		}
		if !failed {
			t.Fatal("failed worker result must preserve dispatch error")
		}
		if second := <-done; len(second) != 2 {
			t.Fatalf("want two concurrent batches, got %#v", second)
		}
	case <-time.After(time.Second):
		t.Fatal("dispatch batch did not finish")
	}
}
