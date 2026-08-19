package infrastructure

import "testing"

func TestCompactNonEmptyDoesNotMutateInput(t *testing.T) {
	input := []string{"a", "", "b"}
	before := append([]string(nil), input...)
	got := CompactNonEmpty(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 values, got %d", len(got))
	}
	for i := range before {
		if before[i] != input[i] {
			t.Fatalf("input was mutated at %d", i)
		}
	}
}
