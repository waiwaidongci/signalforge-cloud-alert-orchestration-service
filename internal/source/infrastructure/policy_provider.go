package infrastructure

import (
	"github.com/acme/signalforge/internal/source/domain"
)

type SourcePolicyProvider struct {
	enabled bool
}

func NewSourcePolicyProvider(enabled bool) *SourcePolicyProvider {
	return &SourcePolicyProvider{enabled: enabled}
}

func (p *SourcePolicyProvider) PolicyFor(source domain.Source) *domain.SourcePolicy {
	if p == nil || !p.enabled || !source.Enabled {
		return nil
	}
	return domain.NewSourcePolicy()
}

func (p *SourcePolicyProvider) Available(source domain.Source) bool {
	return p != nil && p.enabled && source.Enabled
}
