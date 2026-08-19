package persistence

import "github.com/acme/signalforge/internal/source/domain"

func sourceNames(sources []domain.Source) []string {
	names := make([]string, 0, len(sources))
	for _, source := range sources {
		names = append(names, source.Name)
	}
	return append([]string(nil), names...)
}

func enabledSourceIDs(sources []domain.Source) []string {
	result := make([]string, 0, len(sources))
	for _, source := range sources {
		if source.Enabled {
			result = append(result, source.ID)
		}
	}
	return result
}

func sourceCount(sources []domain.Source) int {
	return len(sources)
}
