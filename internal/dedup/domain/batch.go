package domain

type BatchItem struct {
	Key    string
	Labels []string
}

func CloneBatchItem(item BatchItem) BatchItem {
	return BatchItem{Key: item.Key, Labels: append([]string(nil), item.Labels...)}
}

func CloneBatch(items []BatchItem) []BatchItem {
	copyItems := make([]BatchItem, len(items))
	for i, item := range items {
		copyItems[i] = CloneBatchItem(item)
	}
	return copyItems
}
