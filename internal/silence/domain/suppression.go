package domain

import (
	"time"

	"github.com/acme/signalforge/internal/shared/matcher"
	"github.com/acme/signalforge/internal/shared/severity"
)

type SuppressionRule struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	SourceMatcher matcher.Selector  `json:"source_matcher"`
	HighSeverity  severity.Severity `json:"high_severity"`
	LowSeverity   severity.Severity `json:"low_severity"`
	Enabled       bool              `json:"enabled"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}
