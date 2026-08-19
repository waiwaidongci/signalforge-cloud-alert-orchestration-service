package infrastructure

import (
	"reflect"
	"testing"

	"github.com/acme/signalforge/internal/silence/domain"
)

func TestCloneRulesDoesNotMutate(t *testing.T) {
	input := []domain.SuppressionRule{{Name: "one"}, {Name: "two"}}
	before := append([]domain.SuppressionRule(nil), input...)
	got := cloneRules(input)
	got[0].Name = "mutated"
	if !reflect.DeepEqual(before, input) {
		t.Fatalf("input was mutated: before=%+v after=%+v", before, input)
	}
}
