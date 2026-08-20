package infrastructure

import "github.com/acme/signalforge/internal/incident/domain"

func ProjectState(state domain.RetryState) string {
	if state == domain.RetryResolved {
		return "completed"
	}
	return string(state)
}
