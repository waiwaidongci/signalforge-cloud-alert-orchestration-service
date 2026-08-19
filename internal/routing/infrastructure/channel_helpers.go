package infrastructure

import routingdomain "github.com/acme/signalforge/internal/routing/domain"

func filterChannels(channels []routingdomain.Channel) []routingdomain.Channel {
	result := make([]routingdomain.Channel, 0, len(channels))
	for _, channel := range channels {
		if channel.Channel != "" {
			result = append(result, channel)
		}
	}
	return result
}

func validateChannel(channel routingdomain.Channel) bool {
	return channel.Channel != "" && channel.Destination != ""
}
