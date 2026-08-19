package domain

import "time"

type Kind string

const (
	KindPrometheus Kind = "prometheus"
	KindLoki       Kind = "loki"
	KindProbe      Kind = "probe"
	KindCustom     Kind = "custom"
)

type Source struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Kind         Kind              `json:"kind"`
	APIKey       string            `json:"api_key,omitempty"`
	Enabled      bool              `json:"enabled"`
	RateLimit    int               `json:"rate_limit"`
	FieldMapping map[string]string `json:"field_mapping"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func (s Source) Mapping(target string) string {
	if s.FieldMapping == nil {
		s.FieldMapping = map[string]string{}
	}
	if value := s.FieldMapping[target]; value != "" {
		return value
	}
	return target
}
