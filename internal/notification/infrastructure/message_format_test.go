package infrastructure

import "testing"

func TestBuildFailureMessageUnknownChannel(t *testing.T) {
	if got := BuildFailureMessage("missing", "dest"); got != "未知通道" {
		t.Fatalf("expected unknown channel label, got %q", got)
	}
}
