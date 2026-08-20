package db

import (
	"errors"
	"testing"
)

type fakeItemTx struct{ commits, rollbacks int }

func (t *fakeItemTx) Commit() error   { t.commits++; return nil }
func (t *fakeItemTx) Rollback() error { t.rollbacks++; return nil }
func TestProcessItemsRollsBackFailureWithoutCommit(t *testing.T) {
	a, b := &fakeItemTx{}, &fakeItemTx{}
	want := errors.New("work failed")
	err := ProcessItems([]ItemTransaction{a, b}, func(i int) error {
		if i == 1 {
			return want
		}
		return nil
	})
	if !errors.Is(err, want) || a.commits != 1 || b.commits != 0 || b.rollbacks != 1 {
		t.Fatalf("unexpected lifecycle: err=%v a=%#v b=%#v", err, a, b)
	}
}
