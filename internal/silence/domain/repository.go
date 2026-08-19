package domain

import (
	"context"
	"time"
)

type Repository interface {
	CreateSilence(ctx context.Context, silence Silence) error
	UpdateSilence(ctx context.Context, silence Silence) error
	FindSilenceByID(ctx context.Context, id string) (Silence, error)
	ListSilences(ctx context.Context, activeOnly bool, limit, offset int) ([]Silence, int, error)
	DeleteSilence(ctx context.Context, id string) error
	ActiveSilences(ctx context.Context, at time.Time) ([]Silence, error)

	CreateSuppression(ctx context.Context, rule SuppressionRule) error
	UpdateSuppression(ctx context.Context, rule SuppressionRule) error
	FindSuppressionByID(ctx context.Context, id string) (SuppressionRule, error)
	ListSuppressions(ctx context.Context, enabledOnly bool, limit, offset int) ([]SuppressionRule, int, error)
	DeleteSuppression(ctx context.Context, id string) error
	ActiveSuppressions(ctx context.Context) ([]SuppressionRule, error)
}
