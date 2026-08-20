package domain

import "maps"

type IndexEntry struct {
	Fingerprint string
	Labels      map[string]string
}

func CloneIndexEntry(entry IndexEntry) IndexEntry {
	return IndexEntry{Fingerprint: entry.Fingerprint, Labels: maps.Clone(entry.Labels)}
}
