package infrastructure

import "github.com/acme/signalforge/internal/incident/domain"

func ProjectState(state domain.RetryState) string { return string(state) }
