package domain

type IndexEntry struct {
	Fingerprint string
	Labels      map[string]string
}

func CloneIndexEntry(entry IndexEntry) IndexEntry {
	cloned := IndexEntry{Fingerprint: entry.Fingerprint}
	if entry.Labels != nil {
		labels := make(map[string]string, len(entry.Labels))
		for k, v := range entry.Labels {
			labels[k] = v
		}
		cloned.Labels = labels
	}
	return cloned
}
