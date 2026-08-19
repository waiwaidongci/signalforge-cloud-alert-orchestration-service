package infrastructure

import "github.com/acme/signalforge/internal/incident/domain"

func EventTypeRank(eventType domain.EventType) int {
	switch eventType {
	case domain.EventReceived:
		return 0
	case domain.EventDeduped:
		return 1
	case domain.EventAggregated:
		return 2
	case domain.EventNotified:
		return 3
	default:
		return 9
	}
}

func EventTypeName(eventType domain.EventType) string {
	return string(eventType)
}

func EventIDs(events []domain.TimelineEvent) []string {
	ids := make([]string, 0, len(events))
	for _, event := range events {
		ids = append(ids, event.ID)
	}
	return append([]string(nil), ids...)
}
