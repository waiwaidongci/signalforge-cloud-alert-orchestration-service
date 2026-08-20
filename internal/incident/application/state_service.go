package application

import "github.com/acme/signalforge/internal/incident/domain"

func ResolveRetry(state domain.RetryState) domain.RetryState { return domain.AdvanceState(state, true) }
func ActiveIncidents(states []domain.RetryState) int {
	n := 0
	for _, state := range states {
		if domain.IsActiveHistory(state) {
			n++
		}
	}
	return n
}
