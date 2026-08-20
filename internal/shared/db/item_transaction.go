package db

type ItemTransaction interface {
	Commit() error
	Rollback() error
}

func CommitItem(tx ItemTransaction, work func() error) (err error) {
	if err = work(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
