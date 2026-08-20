package http

import "github.com/acme/signalforge/internal/incident/domain"

func detachTimeline(events []domain.TimelineEvent) []domain.TimelineEvent { return events }
