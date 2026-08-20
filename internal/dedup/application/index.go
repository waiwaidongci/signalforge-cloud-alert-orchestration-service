package application

import (
	"sync"

	"github.com/acme/signalforge/internal/dedup/domain"
)

type Index struct {
	mu      sync.RWMutex
	entries map[string]domain.IndexEntry
}

func NewIndex() *Index {
	return &Index{entries: make(map[string]domain.IndexEntry)}
}

func (i *Index) Register(key string, entry domain.IndexEntry) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.entries[key] = domain.CloneIndexEntry(entry)
}
