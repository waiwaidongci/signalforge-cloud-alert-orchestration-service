package infrastructure

import (
	"context"
	"time"

	"github.com/acme/signalforge/internal/incident/domain"
)

type InMemoryTimeline struct {
	events []domain.TimelineEvent
}

func NewInMemoryTimeline() *InMemoryTimeline {
	return &InMemoryTimeline{}
}

func (m *InMemoryTimeline) Append(_ context.Context, event domain.TimelineEvent) error {
	event.CreatedAt = time.Now().UTC()
	m.events = append(m.events, event)
	return nil
}

func (m *InMemoryTimeline) List(_ context.Context, incidentID string, limit, offset int) ([]domain.TimelineEvent, int, error) {
	var result []domain.TimelineEvent
	for _, event := range m.events {
		if event.IncidentID == incidentID {
			result = append(result, event)
		}
	}
	total := len(result)
	if offset >= total {
		return []domain.TimelineEvent{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	page := make([]domain.TimelineEvent, end-offset)
	for i, event := range result[offset:end] {
		page[i] = event.Clone()
	}
	return page, total, nil
}
