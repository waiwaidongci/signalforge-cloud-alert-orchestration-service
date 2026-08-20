package application

import (
	"testing"

	"github.com/acme/signalforge/internal/source/domain"
	"github.com/acme/signalforge/internal/source/infrastructure"
)

func TestResolvePolicyHandlesDisabledProviderAndZeroValue(t *testing.T) {
	provider := infrastructure.NewSourcePolicyProvider(false)
	if provider.Available(domain.Source{Enabled: true}) {
		t.Fatal("disabled provider must not report available")
	}
	policy := ResolvePolicy(provider, domain.Source{Enabled: true})
	if policy == nil {
		t.Fatal("disabled policy provider must resolve to a usable policy")
	}
	if err := policy.AddRule("cluster"); err != nil {
		t.Fatalf("add rule: %v", err)
	}
	if policy.Allows(map[string]string{}) {
		t.Fatal("resolved policy must retain added required rule")
	}
	if !policy.Allows(map[string]string{"cluster": "prod"}) {
		t.Fatal("policy should allow labels that satisfy its rules")
	}
	if !PolicyReady(policy) {
		t.Fatal("policy with a required label must be ready")
	}
	constructed := domain.NewSourcePolicy()
	if err := constructed.AddRule("region"); err != nil {
		t.Fatalf("constructed policy must own rule map: %v", err)
	}
	var zero domain.SourcePolicy
	if err := zero.AddRule("team"); err != nil {
		t.Fatalf("zero-value policy must initialize its rule map: %v", err)
	}
}
