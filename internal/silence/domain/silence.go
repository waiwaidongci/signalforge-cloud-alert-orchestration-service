package domain

import (
	"time"

	"github.com/acme/signalforge/internal/shared/matcher"
)

type Silence struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Matcher     matcher.Selector `json:"matcher"`
	StartsAt    time.Time        `json:"starts_at"`
	EndsAt      time.Time        `json:"ends_at"`
	CreatedBy   string           `json:"created_by"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func (s Silence) Active(at time.Time) bool {
	return !at.Before(s.StartsAt) && !at.After(s.EndsAt)
}
