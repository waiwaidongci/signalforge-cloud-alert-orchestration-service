package application

import "github.com/acme/signalforge/internal/source/domain"

type PolicyProvider interface {
	PolicyFor(domain.Source) *domain.SourcePolicy
}

func ResolvePolicy(provider PolicyProvider, source domain.Source) *domain.SourcePolicy {
	if provider != nil {
		if policy := provider.PolicyFor(source); policy != nil {
			return policy
		}
	}
	return domain.NewSourcePolicy()
}
