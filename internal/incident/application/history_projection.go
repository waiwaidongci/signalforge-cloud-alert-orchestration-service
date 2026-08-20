package application

import "github.com/acme/signalforge/internal/incident/domain"

func HistoryLabel(state domain.RetryState) string {
	if state == domain.RetryResolved {
		return "completed"
	}
	return "active"
}
