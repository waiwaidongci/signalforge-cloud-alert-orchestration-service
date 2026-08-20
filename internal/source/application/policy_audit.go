package application

import "github.com/acme/signalforge/internal/source/domain"

func PolicyReady(policy *domain.SourcePolicy) bool {
	if policy == nil {
		return false
	}
	return !policy.Allows(map[string]string{})
}
