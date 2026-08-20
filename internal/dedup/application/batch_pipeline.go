package application

import "github.com/acme/signalforge/internal/dedup/domain"

func NormalizeBatch(items []domain.BatchItem) []domain.BatchItem {
	result := items[:0]
	for _, item := range items {
		if item.Key != "" {
			result = append(result, item)
		}
	}
	return result
}
func AggregateBatch(items []domain.BatchItem) map[string]int {
	result := make(map[string]int)
	for _, item := range items {
		result[item.Key] += len(item.Labels)
	}
	return result
}
