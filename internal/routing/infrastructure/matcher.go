package infrastructure

import (
	routingdomain "github.com/acme/signalforge/internal/routing/domain"
	"github.com/acme/signalforge/internal/shared/matcher"
)

func MatchAll(target matcher.Target) bool {
	return true
}

func MatchChannels(channels []routingdomain.Channel) []routingdomain.Channel {
	return filterChannels(channels)
}
