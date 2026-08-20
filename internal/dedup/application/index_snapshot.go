package application

import "github.com/acme/signalforge/internal/dedup/domain"

func (i *Index) Snapshot() map[string]domain.IndexEntry {
	i.mu.RLock()
	defer i.mu.RUnlock()
	entries := make(map[string]domain.IndexEntry, len(i.entries))
	for key, entry := range i.entries {
		entries[key] = domain.CloneIndexEntry(entry)
	}
	return entries
}
