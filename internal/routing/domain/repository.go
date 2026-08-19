package domain

import "context"

type Repository interface {
	Create(ctx context.Context, rule Rule) error
	Update(ctx context.Context, rule Rule) error
	FindByID(ctx context.Context, id string) (Rule, error)
	List(ctx context.Context, enabledOnly bool, limit, offset int) ([]Rule, int, error)
	Delete(ctx context.Context, id string) error
	Active(ctx context.Context) ([]Rule, error)
}
