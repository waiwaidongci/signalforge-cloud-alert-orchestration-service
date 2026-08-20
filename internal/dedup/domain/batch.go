package domain

type BatchItem struct {
	Key    string
	Labels []string
}

func CloneBatchItem(item BatchItem) BatchItem { return item }

func CloneBatch(items []BatchItem) []BatchItem { return items }
