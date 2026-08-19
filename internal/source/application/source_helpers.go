package application

import "github.com/acme/signalforge/internal/source/domain"

func filterEnabled(sources []domain.Source) []domain.Source {
	result := make([]domain.Source, 0, len(sources))
	for _, source := range sources {
		if source.Enabled {
			result = append(result, source)
		}
	}
	return result
}

func SourceFilterSummary(sources []domain.Source) string {
	if len(sources) == 0 {
		return "empty"
	}
	return "filtered-sources"
}

func FilteredCount(sources []domain.Source) int {
	return len(sources)
}
