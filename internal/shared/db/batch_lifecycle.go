package db

func ProcessItems(items []ItemTransaction, work func(int) error) error {
	for i, tx := range items {
		defer tx.Rollback()
		if err := CommitItem(tx, func() error { return work(i) }); err != nil {
			return err
		}
	}
	return nil
}
