package infrastructure

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestFanOutCompletesWithoutErrors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	start := make(chan struct{})
	var ready sync.WaitGroup
	for i := 0; i < 2; i++ {
		ready.Add(1)
		go func() {
			defer ready.Done()
			<-start
		}()
	}
	close(start)
	ready.Wait()
	done := make(chan error, 1)
	go func() {
		done <- FanOut(ctx, []func() error{
			func() error { return nil },
			func() error { return nil },
		})
	}()
	select {
	case <-ctx.Done():
		t.Fatal("FanOut hung waiting for channel close")
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	}
}
