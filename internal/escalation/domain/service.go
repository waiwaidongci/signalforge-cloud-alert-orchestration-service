package domain

import "context"

type UseCase interface {
	Create(ctx context.Context, policy Policy) (Policy, error)
	Update(ctx context.Context, policy Policy) (Policy, error)
	Get(ctx context.Context, id string) (Policy, error)
	List(ctx context.Context, enabledOnly bool, limit, offset int) ([]Policy, int, error)
	Delete(ctx context.Context, id string) error
}
