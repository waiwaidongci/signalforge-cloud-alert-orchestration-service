package domain

import (
	"context"
	"errors"
)

var (
	ErrUnknownChannel      = errors.New("unknown notification channel")
	ErrDestinationRequired = errors.New("notification destination is required")
)

type Channel interface {
	Name() string
	Send(ctx context.Context, destination string, payload map[string]any) error
}

type Dispatcher interface {
	Send(ctx context.Context, notification Notification) error
}
