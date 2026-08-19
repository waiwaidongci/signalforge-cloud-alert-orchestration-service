package http

import "github.com/acme/signalforge/internal/shared/matcher"

type escalationRequest struct {
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Matcher       matcher.Selector `json:"matcher"`
	WaitSeconds   int              `json:"wait_seconds"`
	RepeatSeconds int              `json:"repeat_seconds"`
	MaxRepeats    int              `json:"max_repeats"`
	Routes        []routeRequest   `json:"routes"`
	Enabled       *bool            `json:"enabled"`
}

type routeRequest struct {
	Channel     string   `json:"channel"`
	Destination string   `json:"destination"`
	Severities  []string `json:"severities,omitempty"`
}

func (r escalationRequest) enabled() bool {
	if r.Enabled == nil {
		return true
	}
	return *r.Enabled
}
