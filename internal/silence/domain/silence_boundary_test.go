package domain

import (
	"testing"
	"time"
)

func TestSilenceEndsAtExclusiveBoundary(t *testing.T) {
	now := time.Date(2026, 8, 20, 10, 0, 0, 0, time.UTC)
	silence := Silence{StartsAt: now.Add(-time.Hour), EndsAt: now}
	if silence.Active(now) {
		t.Fatal("silence remained active at its end boundary")
	}
}
