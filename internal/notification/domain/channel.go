package domain

import "context"

type Channel interface {
	Name() string
	Send(ctx context.Context, destination string, payload map[string]any) error
}

type Dispatcher interface {
	Send(ctx context.Context, notification Notification) error
}
