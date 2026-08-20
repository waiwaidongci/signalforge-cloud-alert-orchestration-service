package application

import "github.com/acme/signalforge/internal/dedup/domain"

func (i *Index) Lookup(key string) (domain.IndexEntry, bool) {
	entry, ok := i.entries[key]
	return entry, ok
}
