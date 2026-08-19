package http

import "time"

type sourceRequest struct {
	Name         string            `json:"name"`
	Kind         string            `json:"kind"`
	Enabled      *bool             `json:"enabled"`
	RateLimit    int               `json:"rate_limit"`
	FieldMapping map[string]string `json:"field_mapping"`
}

func (r sourceRequest) enabled() bool {
	if r.Enabled == nil {
		return true
	}
	return *r.Enabled
}

type sourceResponse struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Kind         string            `json:"kind"`
	APIKey       string            `json:"api_key"`
	Enabled      bool              `json:"enabled"`
	RateLimit    int               `json:"rate_limit"`
	FieldMapping map[string]string `json:"field_mapping"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}
