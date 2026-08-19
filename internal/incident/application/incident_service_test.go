package application

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/acme/signalforge/internal/incident/domain"
	"github.com/acme/signalforge/internal/shared/clock"
)

type fakeIncidentTimelineForFanout struct {
	err error
}

func (f fakeIncidentTimelineForFanout) Append(context.Context, domain.TimelineEvent) error { return f.err }
func (f fakeIncidentTimelineForFanout) List(context.Context, string, int, int) ([]domain.TimelineEvent, int, error) {
	return nil, 0, nil
}

func TestAppendCompletesWithoutError(t *testing.T) {
	service := NewService(nil, fakeIncidentTimelineForFanout{}, clock.FixedClock{Time: time.Unix(0, 0).UTC()})
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
	go func() { done <- service.Append(ctx, domain.TimelineEvent{}) }()
	select {
	case <-ctx.Done():
		t.Fatal("Append hung waiting for error channel")
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	}
}
