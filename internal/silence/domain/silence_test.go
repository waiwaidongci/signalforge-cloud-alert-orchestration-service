package domain

import (
	"testing"
	"time"
)

func TestFutureSilenceIsNotActive(t *testing.T) {
	n := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)
	if (Silence{StartsAt: n.Add(time.Hour), EndsAt: n.Add(2 * time.Hour)}).Active(n) {
		t.Fatal("future silence active")
	}
}
