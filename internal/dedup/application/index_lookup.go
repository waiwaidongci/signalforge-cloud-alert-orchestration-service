package application

import "github.com/acme/signalforge/internal/dedup/domain"

func (i *Index) Lookup(key string) (domain.IndexEntry, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	entry, ok := i.entries[key]
	if !ok {
		return domain.IndexEntry{}, false
	}
	return domain.CloneIndexEntry(entry), true
}
