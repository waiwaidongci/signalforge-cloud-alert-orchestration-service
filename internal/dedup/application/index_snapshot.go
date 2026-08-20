package application

import "github.com/acme/signalforge/internal/dedup/domain"

func (i *Index) Snapshot() map[string]domain.IndexEntry {
	i.mu.RLock()
	defer i.mu.RUnlock()
	out := make(map[string]domain.IndexEntry, len(i.entries))
	for k, v := range i.entries {
		out[k] = domain.CloneIndexEntry(v)
	}
	return out
}
