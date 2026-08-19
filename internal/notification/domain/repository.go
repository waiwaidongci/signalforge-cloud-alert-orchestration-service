package domain

import "context"

type Repository interface {
	Create(ctx context.Context, notification Notification) error
	Update(ctx context.Context, notification Notification) error
	FindByID(ctx context.Context, id string) (Notification, error)
	ListByIncident(ctx context.Context, incidentID string, limit, offset int) ([]Notification, int, error)
	List(ctx context.Context, channel string, status Status, limit, offset int) ([]Notification, int, error)
}
