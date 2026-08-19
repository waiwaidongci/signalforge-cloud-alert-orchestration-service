package infrastructure

import (
	"sync"
	"testing"
)

func TestFanoutStateBusyWhenWorkExists(t *testing.T) {
	funcs := []func() error{func() error { return nil }}
	start := make(chan struct{})
	var ready sync.WaitGroup
	for i := 0; i < 2; i++ {
		ready.Add(1)
		go func() {
			defer ready.Done()
			<-start
			_ = FanoutState(funcs)
		}()
	}
	close(start)
	ready.Wait()
	if got := FanoutState(funcs); got != "busy" {
		t.Fatalf("expected busy, got %q", got)
	}
}
