package domain

import "context"

type Repository interface {
	Create(ctx context.Context, source Source) error
	Update(ctx context.Context, source Source) error
	FindByID(ctx context.Context, id string) (Source, error)
	FindByAPIKey(ctx context.Context, apiKey string) (Source, error)
	List(ctx context.Context, enabled *bool, limit, offset int) ([]Source, int, error)
	Delete(ctx context.Context, id string) error
}
