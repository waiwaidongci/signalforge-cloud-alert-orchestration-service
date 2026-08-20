package application

import "github.com/acme/signalforge/internal/incident/domain"

func HistoryLabel(state domain.RetryState) string { return "active" }
