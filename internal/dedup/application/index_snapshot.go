package application

import "github.com/acme/signalforge/internal/dedup/domain"

func (i *Index) Snapshot() map[string]domain.IndexEntry {
	return i.entries
}
