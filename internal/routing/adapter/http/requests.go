package http

import "github.com/acme/signalforge/internal/shared/matcher"

type routingRequest struct {
	Name     string           `json:"name"`
	Priority int              `json:"priority"`
	Match    matcher.Selector `json:"match"`
	Channels []channelRequest `json:"channels"`
	Enabled  *bool            `json:"enabled"`
}

type channelRequest struct {
	Channel     string `json:"channel"`
	Destination string `json:"destination"`
}

func (r routingRequest) enabled() bool {
	if r.Enabled == nil {
		return true
	}
	return *r.Enabled
}
