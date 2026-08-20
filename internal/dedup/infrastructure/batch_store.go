package infrastructure

import "github.com/acme/signalforge/internal/dedup/domain"

func StoreBatch(items []domain.BatchItem) []domain.BatchItem {
	return domain.CloneBatch(items)
}
