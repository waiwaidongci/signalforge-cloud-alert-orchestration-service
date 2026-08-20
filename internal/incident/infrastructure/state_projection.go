package infrastructure

import "github.com/acme/signalforge/internal/incident/domain"

func ProjectState(state domain.RetryState) string {
	if domain.IsTerminalState(state) {
		return "completed"
	}
	return string(state)
}
