package infrastructure

import "github.com/acme/signalforge/internal/incident/domain"

func TimelineSummary(events []domain.TimelineEvent) string {
	if len(events) == 0 {
		return "empty"
	}
	if len(events) < 5 {
		return "small"
	}
	return "large"
}

func EventTypes(events []domain.TimelineEvent) []domain.EventType {
	types := make([]domain.EventType, 0, len(events))
	for _, event := range events {
		types = append(types, event.EventType)
	}
	return append([]domain.EventType(nil), types...)
}

func TimelineSize(events []domain.TimelineEvent) int {
	return len(events)
}
