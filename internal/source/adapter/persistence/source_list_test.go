package persistence

import (
	"reflect"
	"testing"

	"github.com/acme/signalforge/internal/source/domain"
)

func TestCompactSourcesDoesNotMutateInput(t *testing.T) {
	input := []domain.Source{
		{Name: "one", Enabled: true},
		{Name: "two", Enabled: false},
		{Name: "three", Enabled: true},
	}
	before := append([]domain.Source(nil), input...)
	got := compactSources(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(got))
	}
	for i := range before {
		if !reflect.DeepEqual(before[i], input[i]) {
			t.Fatalf("input was mutated at %d", i)
		}
	}
}
