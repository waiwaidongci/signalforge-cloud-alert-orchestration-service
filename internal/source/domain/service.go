package domain

import "context"

type UseCase interface {
	Create(ctx context.Context, source Source) (Source, error)
	Update(ctx context.Context, source Source) (Source, error)
	Get(ctx context.Context, id string) (Source, error)
	List(ctx context.Context, enabled *bool, limit, offset int) ([]Source, int, error)
	Delete(ctx context.Context, id string) error
}
