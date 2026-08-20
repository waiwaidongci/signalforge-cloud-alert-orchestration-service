package db

import "errors"

type ItemTransaction interface {
	Commit() error
	Rollback() error
}

func RollbackItem(tx ItemTransaction, err error) error {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}
	return err
}

func CommitItem(tx ItemTransaction, work func() error) error {
	if err := work(); err != nil {
		return RollbackItem(tx, err)
	}
	return tx.Commit()
}
