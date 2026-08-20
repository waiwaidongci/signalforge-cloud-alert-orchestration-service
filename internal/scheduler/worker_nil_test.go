package scheduler

import (
	"context"
	"testing"
)

func TestRecoveryWorkerSkipsMissingDependencies(t *testing.T) {
	var runner *Runner
	defer func() {
		if recover() != nil {
			t.Fatal("recovery worker panicked with missing dependencies")
		}
	}()
	runner.runRecovery(context.Background())
}

func TestSilenceExpirySkipsMissingDependencies(t *testing.T) {
	var runner *Runner
	defer func() {
		if recover() != nil {
			t.Fatal("silence worker panicked with missing dependencies")
		}
	}()
	runner.runSilenceExpiry(context.Background())
}

func TestSchedulerLoggerHandlesNil(t *testing.T) {
	if schedulerLogger(nil) == nil {
		t.Fatal("scheduler logger must fall back to a usable logger")
	}
}
