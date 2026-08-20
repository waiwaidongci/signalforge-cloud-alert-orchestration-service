package application

import "github.com/acme/signalforge/internal/source/domain"

func PolicyReady(policy *domain.SourcePolicy) bool { return policy != nil && len(policy.RequiredLabels()) > 0 }
