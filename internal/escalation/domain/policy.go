package domain

import (
	"time"

	"github.com/acme/signalforge/internal/shared/matcher"
)

type Route struct {
	Channel     string   `json:"channel"`
	Destination string   `json:"destination"`
	Severities  []string `json:"severities,omitempty"`
}

type Policy struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Matcher       matcher.Selector `json:"matcher"`
	WaitSeconds   int              `json:"wait_seconds"`
	RepeatSeconds int              `json:"repeat_seconds"`
	MaxRepeats    int              `json:"max_repeats"`
	Routes        []Route          `json:"routes"`
	Enabled       bool             `json:"enabled"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

func (p Policy) NextDelay(repeat int) time.Duration {
	if repeat == 0 {
		return time.Duration(p.WaitSeconds) * time.Second
	}
	return time.Duration(p.RepeatSeconds) * time.Second
}

func (p Policy) AllowsRepeat(repeat int) bool {
	return repeat < p.MaxRepeats
}
