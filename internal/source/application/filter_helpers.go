package application

import "github.com/acme/signalforge/internal/source/domain"

func SelectEnabled(sources []domain.Source) []domain.Source {
	result := sources[:0]
	for _, source := range sources {
		if source.Enabled {
			result = append(result, source)
		}
	}
	return result
}

func SourceKey(source domain.Source) string {
	return source.ID
}

func EnabledCount(sources []domain.Source) int {
	count := 0
	for _, source := range sources {
		if source.Enabled {
			count++
		}
	}
	return count
}

func FilterSummary(sources []domain.Source) string {
	if len(sources) == 0 {
		return "empty"
	}
	return "enabled"
}
