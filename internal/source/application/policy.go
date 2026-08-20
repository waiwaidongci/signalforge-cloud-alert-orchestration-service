package application

import "github.com/acme/signalforge/internal/source/domain"

type PolicyProvider interface {
	PolicyFor(domain.Source) *domain.SourcePolicy
}

func ResolvePolicy(provider PolicyProvider, source domain.Source) *domain.SourcePolicy {
	if provider == nil {
		return domain.NewSourcePolicy()
	}
	policy := provider.PolicyFor(source)
	if policy == nil {
		return domain.NewSourcePolicy()
	}
	return policy
}
