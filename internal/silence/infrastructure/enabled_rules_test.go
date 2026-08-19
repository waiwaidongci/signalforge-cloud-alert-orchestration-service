package infrastructure

import (
	"reflect"
	"testing"

	"github.com/acme/signalforge/internal/silence/domain"
)

func TestEnabledRulesDoesNotMutateInput(t *testing.T) {
	input := []domain.SuppressionRule{{Name: "one", Enabled: false}, {Name: "two", Enabled: true}}
	before := append([]domain.SuppressionRule(nil), input...)
	got := EnabledRules(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 enabled rule, got %d", len(got))
	}
	if !reflect.DeepEqual(before, input) {
		t.Fatalf("input was mutated: before=%+v after=%+v", before, input)
	}
}
