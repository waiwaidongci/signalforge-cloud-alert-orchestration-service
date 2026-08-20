package domain

type BatchItem struct {
	Key    string
	Labels []string
}

func CloneBatchItem(item BatchItem) BatchItem {
	cloned := item
	if item.Labels != nil {
		cloned.Labels = append([]string(nil), item.Labels...)
	}
	return cloned
}

func CloneBatch(items []BatchItem) []BatchItem {
	if items == nil {
		return nil
	}
	cloned := make([]BatchItem, len(items))
	for i, item := range items {
		cloned[i] = CloneBatchItem(item)
	}
	return cloned
}
