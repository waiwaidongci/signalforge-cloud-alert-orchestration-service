package db

import "errors"

type ItemTransaction interface {
	Commit() error
	Rollback() error
}

func CommitItem(tx ItemTransaction, work func() error) (err error) {
	defer func() { err = tx.Commit() }()
	if err := work(); err != nil {
		return err
	}
	return nil
}

var _ = errors.New
