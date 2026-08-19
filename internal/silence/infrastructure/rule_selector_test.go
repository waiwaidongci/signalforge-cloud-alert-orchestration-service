package infrastructure

import (
	"testing"

	"github.com/acme/signalforge/internal/silence/domain"
)

func TestSelectedRuleNameUsesLastRule(t *testing.T) {
	rules := []domain.SuppressionRule{{Name: "first"}, {Name: "last"}}
	if got := SelectedRuleName(rules); got != "last" {
		t.Fatalf("expected last, got %q", got)
	}
}
