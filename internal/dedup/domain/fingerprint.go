package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/acme/signalforge/internal/shared/severity"
)

type Input struct {
	SourceID string
	Resource string
	Severity severity.Severity
	Title    string
	Labels   map[string]string
}

func Fingerprint(input Input) string {
	labelKeys := make([]string, 0, len(input.Labels))
	for key := range input.Labels {
		labelKeys = append(labelKeys, key)
	}
	sort.Strings(labelKeys)

	var builder strings.Builder
	builder.WriteString(input.SourceID)
	builder.WriteByte(0)
	builder.WriteString(input.Resource)
	builder.WriteByte(0)
	builder.WriteString(input.Severity.String())
	builder.WriteByte(0)
	builder.WriteString(input.Title)
	for _, key := range labelKeys {
		builder.WriteByte(0)
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(input.Labels[key])
	}
	sum := sha256.Sum256([]byte(builder.String()))
	return hex.EncodeToString(sum[:])
}
