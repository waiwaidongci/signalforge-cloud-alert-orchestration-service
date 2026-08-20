package db

func ProcessItems(items []ItemTransaction, work func(int) error) error {
	for i, tx := range items {
		if err := CommitItem(tx, func() error { return work(i) }); err != nil {
			return err
		}
	}
	return nil
}
