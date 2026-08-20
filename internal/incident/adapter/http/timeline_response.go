package http

import "github.com/acme/signalforge/internal/incident/domain"

func detachTimeline(events []domain.TimelineEvent) []domain.TimelineEvent {
	result := make([]domain.TimelineEvent, len(events))
	for i, event := range events {
		result[i] = event.Clone()
	}
	return result
}
