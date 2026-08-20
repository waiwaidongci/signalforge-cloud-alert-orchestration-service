package matcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/acme/signalforge/internal/shared/severity"
)

var ErrInvalidSelector = errors.New("invalid selector")

type Selector struct {
	SourceIDs  []string          `json:"source_ids,omitempty"`
	Services   []string          `json:"services,omitempty"`
	Envs       []string          `json:"environments,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	Severities []string          `json:"severities,omitempty"`
}

type Target struct {
	SourceID string            `json:"source_id"`
	Service  string            `json:"service"`
	Env      string            `json:"environment"`
	Labels   map[string]string `json:"labels"`
	Severity severity.Severity `json:"severity"`
}

func (s Selector) Matches(target Target) bool {
	if len(s.SourceIDs) > 0 && !contains(s.SourceIDs, target.SourceID) {
		return false
	}
	if len(s.Services) > 0 && !contains(s.Services, target.Service) {
		return false
	}
	if len(s.Envs) > 0 && !contains(s.Envs, target.Env) {
		return false
	}
	if len(s.Severities) > 0 && !contains(s.Severities, target.Severity.String()) {
		return false
	}
	for key, value := range s.Labels {
		if target.Labels[key] != value {
			return false
		}
	}
	return true
}

func EncodeSelector(selector Selector) string {
	raw, _ := json.Marshal(selector)
	return string(raw)
}

func DecodeSelector(raw string) (Selector, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Selector{}, nil
	}
	if trimmed == "null" {
		return Selector{}, fmt.Errorf("decode selector: %w", ErrInvalidSelector)
	}
	var selector Selector
	if err := json.Unmarshal([]byte(raw), &selector); err != nil {
		return Selector{}, fmt.Errorf("decode selector: %w", errors.Join(ErrInvalidSelector, err))
	}
	return selector, nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
