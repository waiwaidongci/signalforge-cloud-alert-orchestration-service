package http

import (
	"time"

	"github.com/acme/signalforge/internal/shared/matcher"
)

type silenceRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Matcher     matcher.Selector `json:"matcher"`
	StartsAt    time.Time        `json:"starts_at"`
	EndsAt      time.Time        `json:"ends_at"`
	CreatedBy   string           `json:"created_by"`
}

type suppressionRequest struct {
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	SourceMatcher matcher.Selector `json:"source_matcher"`
	HighSeverity  string           `json:"high_severity"`
	LowSeverity   string           `json:"low_severity"`
	Enabled       *bool            `json:"enabled"`
}

func (r suppressionRequest) enabled() bool {
	if r.Enabled == nil {
		return true
	}
	return *r.Enabled
}
