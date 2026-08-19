package infrastructure

import "testing"

func TestMergeEventsDoesNotMutateLeft(t *testing.T) {
	left := make([]string, 1, 2)
	left[0] = "one"
	right := []string{"two"}
	before := append([]string(nil), left...)
	got := mergeEvents(left, right)
	if len(got) != 2 {
		t.Fatalf("expected 2 events, got %d", len(got))
	}
	if len(left) != len(before) || left[0] != before[0] {
		t.Fatalf("left was mutated: before=%v after=%v", before, left)
	}
	if left[:2][1] != "" {
		t.Fatalf("left backing array was mutated: %v", left[:2])
	}
}
