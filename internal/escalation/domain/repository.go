package domain

import "context"

type Repository interface {
	Create(ctx context.Context, policy Policy) error
	Update(ctx context.Context, policy Policy) error
	FindByID(ctx context.Context, id string) (Policy, error)
	List(ctx context.Context, enabledOnly bool, limit, offset int) ([]Policy, int, error)
	Delete(ctx context.Context, id string) error
	Active(ctx context.Context) ([]Policy, error)
}
