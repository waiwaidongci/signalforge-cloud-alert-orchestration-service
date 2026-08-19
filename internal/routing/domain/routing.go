package domain

import (
	"time"

	"github.com/acme/signalforge/internal/shared/matcher"
)

type Channel struct {
	Channel     string `json:"channel"`
	Destination string `json:"destination"`
}

type Rule struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Priority  int              `json:"priority"`
	Match     matcher.Selector `json:"match"`
	Channels  []Channel        `json:"channels"`
	Enabled   bool             `json:"enabled"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
