package domain

import "context"

type UseCase interface {
	FindByID(ctx context.Context, id string) (Incident, error)
	List(ctx context.Context, filter ListFilter, limit, offset int) ([]Incident, int, error)
	Acknowledge(ctx context.Context, ids []string, actor string) error
	Close(ctx context.Context, ids []string, actor string) error
	Timeline(ctx context.Context, id string, limit, offset int) ([]TimelineEvent, int, error)
}
