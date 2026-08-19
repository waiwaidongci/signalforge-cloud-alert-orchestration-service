package application

import (
	"reflect"
	"testing"

	"github.com/acme/signalforge/internal/source/domain"
)

func TestSelectEnabledDoesNotMutate(t *testing.T) {
	input := []domain.Source{{Name: "one", Enabled: false}, {Name: "two", Enabled: true}}
	before := append([]domain.Source(nil), input...)
	got := SelectEnabled(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 enabled source, got %d", len(got))
	}
	if !reflect.DeepEqual(before, input) {
		t.Fatalf("input was mutated: before=%+v after=%+v", before, input)
	}
}
