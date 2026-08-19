package persistence

import "github.com/acme/signalforge/internal/source/domain"

func SourceSummary(sources []domain.Source) string {
	if len(sources) == 0 {
		return "none"
	}
	if len(sources) < 4 {
		return "few"
	}
	return "many"
}

func EnabledSourceNames(sources []domain.Source) []string {
	names := make([]string, 0, len(sources))
	for _, source := range sources {
		if source.Enabled {
			names = append(names, source.Name)
		}
	}
	return append([]string(nil), names...)
}

func SourceIDs(sources []domain.Source) []string {
	ids := make([]string, 0, len(sources))
	for _, source := range sources {
		ids = append(ids, source.ID)
	}
	return append([]string(nil), ids...)
}
