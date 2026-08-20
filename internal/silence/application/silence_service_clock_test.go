package application

import (
	"context"
	"testing"
	"time"

	"github.com/acme/signalforge/internal/shared/clock"
	"github.com/acme/signalforge/internal/silence/domain"
)

type clockRecordingRepository struct {
	domain.Repository
	calledAt time.Time
}

func (r *clockRecordingRepository) ListSilences(_ context.Context, _ bool, at time.Time, _, _ int) ([]domain.Silence, int, error) {
	r.calledAt = at
	return nil, 0, nil
}

func TestListSilencesUsesServiceClock(t *testing.T) {
	want := time.Date(2026, 8, 20, 11, 0, 0, 0, time.UTC)
	repository := &clockRecordingRepository{}
	service := NewService(repository, clock.FixedClock{Time: want})

	if _, _, err := service.ListSilences(context.Background(), true, 25, 0); err != nil {
		t.Fatalf("list silences: %v", err)
	}
	if !repository.calledAt.Equal(want) {
		t.Fatalf("repository received %s, want service clock %s", repository.calledAt, want)
	}
}
