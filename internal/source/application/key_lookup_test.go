package application

import (
	"errors"
	"testing"

	"github.com/acme/signalforge/internal/source/domain"
)

func TestResolveKeyPreservesNotFoundChain(t *testing.T) {
	_, err := ResolveKey("missing")
	if !errors.Is(err, domain.ErrAPIKeyNotFound) {
		t.Fatalf("not found sentinel must survive lookup chain: %v", err)
	}
	if got := ClassifyKey(err); got != "not_found" {
		t.Fatalf("expected not_found classification, got %q", got)
	}
}
