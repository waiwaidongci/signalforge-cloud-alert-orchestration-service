package application

import "github.com/acme/signalforge/internal/dedup/domain"

func NormalizeBatch(items []domain.BatchItem) []domain.BatchItem {
	result := make([]domain.BatchItem, 0, len(items))
	for _, item := range items {
		if item.Key != "" {
			result = append(result, domain.CloneBatchItem(item))
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
