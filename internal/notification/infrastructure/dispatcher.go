package infrastructure

import (
	"context"
	"fmt"

	"github.com/acme/signalforge/internal/notification/domain"
)

type Dispatcher struct {
	channels map[string]domain.Channel
}

func NewDispatcher(channels ...domain.Channel) *Dispatcher {
	dispatcher := &Dispatcher{channels: make(map[string]domain.Channel)}
	for _, channel := range channels {
		dispatcher.channels[channel.Name()] = channel
	}
	return dispatcher
}

func (d *Dispatcher) Send(ctx context.Context, notification domain.Notification) error {
	channel := d.channels[notification.Channel]
	if channel == nil {
		return fmt.Errorf("unknown notification channel %q", notification.Channel)
	}
	return channel.Send(ctx, notification.Destination, notification.Payload)
}
