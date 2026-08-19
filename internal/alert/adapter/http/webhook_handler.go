package http

import (
	"log/slog"
	"net/http"
	"strings"

	alertapp "github.com/acme/signalforge/internal/alert/application"
	alertdomain "github.com/acme/signalforge/internal/alert/domain"
	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/httpx"
	sourceapp "github.com/acme/signalforge/internal/source/application"
)

type WebhookHandler struct {
	ingest  *alertapp.IngestService
	sources *sourceapp.Service
	logger  *slog.Logger
}

func NewWebhookHandler(ingest *alertapp.IngestService, sources *sourceapp.Service, logger *slog.Logger) *WebhookHandler {
	return &WebhookHandler{ingest: ingest, sources: sources, logger: logger}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	var payload ingestRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	if payload.Labels == nil {
		payload.Labels = map[string]string{}
	}
	if payload.Annotations == nil {
		payload.Annotations = map[string]string{}
	}
	apiKey := sourceAPIKey(r)
	source, err := h.sources.GetByAPIKey(r.Context(), apiKey)
	if err != nil {
		httpx.WriteError(w, requestID, apperr.New(apperr.KindUnauthorized, "INVALID_SOURCE_KEY", "告警源 API Key 无效"))
		return
	}
	if !source.Enabled {
		httpx.WriteError(w, requestID, apperr.New(apperr.KindForbidden, "SOURCE_DISABLED", "告警源已禁用"))
		return
	}
	result, err := h.ingest.Ingest(r.Context(), alertdomain.IngestInput{
		SourceID:    source.ID,
		ExternalID:  payload.ExternalID,
		Resource:    payload.Resource,
		Severity:    payload.Severity,
		Title:       payload.Title,
		Description: payload.Description,
		Labels:      payload.Labels,
		Annotations: payload.Annotations,
		ReceivedAt:  payload.Timestamp,
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, result, requestID)
}

func sourceAPIKey(r *http.Request) string {
	if value := strings.TrimSpace(r.Header.Get("X-Source-Key")); value != "" {
		return value
	}
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
	}
	return ""
}
