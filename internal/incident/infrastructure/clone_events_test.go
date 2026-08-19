package infrastructure

import (
	"sync"
	"testing"
)

func TestCloneEventsDoesNotMutate(t *testing.T) {
	input := []string{"one", "two"}
	before := append([]string(nil), input...)
	start := make(chan struct{})
	var ready sync.WaitGroup
	for i := 0; i < 2; i++ {
		ready.Add(1)
		go func() {
			defer ready.Done()
			<-start
			_ = cloneEvents(input)
		}()
	}
	close(start)
	ready.Wait()
	got := cloneEvents(input)
	got[0] = "mutated"
	if input[0] != before[0] {
		t.Fatalf("input was mutated: before=%v after=%v", before, input)
	}
}
