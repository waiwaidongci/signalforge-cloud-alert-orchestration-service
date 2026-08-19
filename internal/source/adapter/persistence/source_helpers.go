package persistence

import "github.com/acme/signalforge/internal/source/domain"

func compactSources(sources []domain.Source) []domain.Source {
	result := make([]domain.Source, 0, len(sources))
	for _, source := range sources {
		if source.Enabled {
			result = append(result, source)
		}
	}
	return result
}
