package application

import (
	"sync"
	"testing"
)

func TestCloneEventFuncsDoesNotMutate(t *testing.T) {
	left := make([]func() error, 1, 2)
	left[0] = func() error { return nil }
	before := len(left)
	start := make(chan struct{})
	var ready sync.WaitGroup
	for i := 0; i < 2; i++ {
		ready.Add(1)
		go func() {
			defer ready.Done()
			<-start
			_ = cloneEventFuncs(left)
		}()
	}
	close(start)
	ready.Wait()
	got := cloneEventFuncs(left)
	got = append(got, func() error { return nil })
	if len(left) != before {
		t.Fatalf("left was mutated: before=%d after=%d", before, len(left))
	}
	if left[:2][1] != nil {
		t.Fatalf("left backing array was mutated")
	}
}
