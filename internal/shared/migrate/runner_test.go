package migrate

import (
	"context"
	"errors"
	"testing"
)

type failingDriver struct{ err error }

func (d failingDriver) Apply(context.Context) error { return d.err }
func TestPermanentMigrationErrorPreservesCauseAndSkipsRetry(t *testing.T) {
	err := RunMigrations(context.Background(), failingDriver{err: ErrPermanent})
	if !errors.Is(err, ErrPermanent) || ShouldRetryMigration(err) {
		t.Fatalf("permanent error misclassified: %v", err)
	}
	if MigrationErrorClass(err) != "permanent" {
		t.Fatal("permanent error class lost")
	}
}
