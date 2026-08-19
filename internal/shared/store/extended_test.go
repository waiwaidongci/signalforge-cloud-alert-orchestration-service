package store

import "testing"

func TestExecTraceIncludesContextState(t *testing.T) {
	got := ExecTrace(nil)
	if len(got) != 3 {
		t.Fatalf("expected 3 trace parts, got %v", got)
	}
}
