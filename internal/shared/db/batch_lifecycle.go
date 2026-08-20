package db

func processItem(tx ItemTransaction, work func() error) error { return CommitItem(tx, work) }

func ProcessItems(items []ItemTransaction, work func(int) error) error {
	for i, tx := range items {
		if err := processItem(tx, func() error { return work(i) }); err != nil {
			return err
		}
	}
	return nil
}
