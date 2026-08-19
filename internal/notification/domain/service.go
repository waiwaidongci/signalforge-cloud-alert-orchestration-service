package domain

import "context"

type UseCase interface {
	Notify(ctx context.Context, notification Notification) error
	Get(ctx context.Context, id string) (Notification, error)
	List(ctx context.Context, channel string, status Status, limit, offset int) ([]Notification, int, error)
}
