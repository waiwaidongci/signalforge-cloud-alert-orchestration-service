package persistence

import "testing"

func TestMergeSourcesDoesNotMutateLeft(t *testing.T) {
	left := []string{"one"}
	right := []string{"two"}
	before := append([]string(nil), left...)
	got := mergeSources(left, right)
	if len(got) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(got))
	}
	if len(left) != len(before) || left[0] != before[0] {
		t.Fatalf("left was mutated: before=%v after=%v", before, left)
	}
}
