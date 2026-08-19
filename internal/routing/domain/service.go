package domain

import "context"

type UseCase interface {
	Create(ctx context.Context, rule Rule) (Rule, error)
	Update(ctx context.Context, rule Rule) (Rule, error)
	Get(ctx context.Context, id string) (Rule, error)
	List(ctx context.Context, enabledOnly bool, limit, offset int) ([]Rule, int, error)
	Delete(ctx context.Context, id string) error
}
