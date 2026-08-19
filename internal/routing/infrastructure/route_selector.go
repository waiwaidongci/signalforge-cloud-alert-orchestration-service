package infrastructure

import (
	"strings"

	routingdomain "github.com/acme/signalforge/internal/routing/domain"
)

func ChannelPriority(channel string) int {
	switch channel {
	case "email":
		return 3
	case "webhook":
		return 2
	case "log":
		return 1
	default:
		return 0
	}
}

func SelectChannels(channels []routingdomain.Channel) []routingdomain.Channel {
	result := make([]routingdomain.Channel, 0, len(channels))
	for _, channel := range channels {
		if strings.TrimSpace(channel.Channel) != "" {
			result = append(result, channel)
		}
	}
	return result
}

func ChannelNames(channels []routingdomain.Channel) []string {
	names := make([]string, 0, len(channels))
	for _, channel := range channels {
		names = append(names, channel.Channel)
	}
	return names
}
