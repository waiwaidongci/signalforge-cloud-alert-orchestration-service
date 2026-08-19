package domain

import "context"

type UseCase interface {
	CreateSilence(ctx context.Context, silence Silence) (Silence, error)
	UpdateSilence(ctx context.Context, silence Silence) (Silence, error)
	GetSilence(ctx context.Context, id string) (Silence, error)
	ListSilences(ctx context.Context, activeOnly bool, limit, offset int) ([]Silence, int, error)
	DeleteSilence(ctx context.Context, id string) error

	CreateSuppression(ctx context.Context, rule SuppressionRule) (SuppressionRule, error)
	UpdateSuppression(ctx context.Context, rule SuppressionRule) (SuppressionRule, error)
	GetSuppression(ctx context.Context, id string) (SuppressionRule, error)
	ListSuppressions(ctx context.Context, enabledOnly bool, limit, offset int) ([]SuppressionRule, int, error)
	DeleteSuppression(ctx context.Context, id string) error
}
