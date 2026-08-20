package application

import (
	"github.com/acme/signalforge/internal/incident/domain"
	"github.com/acme/signalforge/internal/incident/infrastructure"
	"testing"
)

func TestResolvedRetryIsProjectedInActiveHistory(t *testing.T) {
	state := ResolveRetry(domain.RetryFailed)
	if state != domain.RetryResolved || ActiveIncidents([]domain.RetryState{state}) != 1 {
		t.Fatalf("retry state not visible: %s", state)
	}
	if infrastructure.ProjectState(state) != "completed" {
		t.Fatal("resolved state projection missing")
	}
	if HistoryLabel(state) != "completed" {
		t.Fatal("history label missing")
	}
}
