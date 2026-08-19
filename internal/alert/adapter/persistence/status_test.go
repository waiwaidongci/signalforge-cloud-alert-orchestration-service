package persistence

import (
	"testing"

	"github.com/acme/signalforge/internal/alert/domain"
)

func TestNormalizeStatusPreservesAcknowledged(t *testing.T) {
	if got := normalizeStatus(domain.StatusAcknowledged); got != domain.StatusAcknowledged {
		t.Fatalf("expected acknowledged, got %q", got)
	}
}
