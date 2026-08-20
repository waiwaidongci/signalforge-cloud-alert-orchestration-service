package persistence

import (
	"errors"
	"testing"

	"github.com/acme/signalforge/internal/escalation/domain"
	"github.com/acme/signalforge/internal/shared/matcher"
)

type policyScanner struct {
	matchJSON  string
	routesJSON string
}

func (s policyScanner) Scan(dest ...any) error {
	values := []any{"p", "n", "", s.matchJSON, 1, 1, 1, s.routesJSON, 1, "", ""}
	for i := range dest {
		switch v := dest[i].(type) {
		case *string:
			*v = values[i].(string)
		case *int:
			*v = values[i].(int)
		}
	}
	return nil
}

func TestScanPolicyRejectsCorruptStoredJSON(t *testing.T) {
	_, err := scanPolicy(policyScanner{matchJSON: `null`, routesJSON: `[]`})
	if !errors.Is(err, domain.ErrInvalidPolicyData) || !errors.Is(err, matcher.ErrInvalidSelector) {
		t.Fatalf("err=%v", err)
	}
}

func TestScanPolicyRejectsCorruptRoutesJSON(t *testing.T) {
	_, err := scanPolicy(policyScanner{matchJSON: `{"services":["api"]}`, routesJSON: `[{"channel":"","destination":""}]`})
	if !errors.Is(err, domain.ErrInvalidRoute) {
		t.Fatalf("err=%v", err)
	}
}
