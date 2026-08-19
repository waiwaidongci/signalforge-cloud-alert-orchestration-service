package http

import (
	"log/slog"
	"net/http"

	"github.com/acme/signalforge/internal/shared/httpx"
	"github.com/acme/signalforge/internal/shared/severity"
	silenceapp "github.com/acme/signalforge/internal/silence/application"
	silencedomain "github.com/acme/signalforge/internal/silence/domain"
)

type Handler struct {
	service *silenceapp.Service
	logger  *slog.Logger
}

func NewHandler(service *silenceapp.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) CreateSilence(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	var payload silenceRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	silence, err := h.service.CreateSilence(r.Context(), silencedomain.Silence{
		Name: payload.Name, Description: payload.Description, Matcher: payload.Matcher, StartsAt: payload.StartsAt, EndsAt: payload.EndsAt, CreatedBy: payload.CreatedBy,
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, silence, requestID)
}

func (h *Handler) UpdateSilence(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	var payload silenceRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	silence, err := h.service.UpdateSilence(r.Context(), silencedomain.Silence{
		ID: id, Name: payload.Name, Description: payload.Description, Matcher: payload.Matcher, StartsAt: payload.StartsAt, EndsAt: payload.EndsAt, CreatedBy: payload.CreatedBy,
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, silence, requestID)
}

func (h *Handler) GetSilence(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	silence, err := h.service.GetSilence(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, silence, requestID)
}

func (h *Handler) ListSilences(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	pagination := httpx.ParsePagination(r.URL.Query())
	activeOnly := r.URL.Query().Get("active") == "true"
	silences, total, err := h.service.ListSilences(r.Context(), activeOnly, pagination.PerPage, pagination.Offset)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteMeta(w, http.StatusOK, silences, httpx.NewPageMeta(pagination, total), requestID)
}

func (h *Handler) DeleteSilence(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	if err := h.service.DeleteSilence(r.Context(), id); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusNoContent, nil, requestID)
}

func (h *Handler) CreateSuppression(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	var payload suppressionRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	rule, err := h.service.CreateSuppression(r.Context(), silencedomain.SuppressionRule{
		Name: payload.Name, Description: payload.Description, SourceMatcher: payload.SourceMatcher,
		HighSeverity: severity.Severity(payload.HighSeverity), LowSeverity: severity.Severity(payload.LowSeverity), Enabled: payload.enabled(),
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, rule, requestID)
}

func (h *Handler) UpdateSuppression(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	var payload suppressionRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	rule, err := h.service.UpdateSuppression(r.Context(), silencedomain.SuppressionRule{
		ID: id, Name: payload.Name, Description: payload.Description, SourceMatcher: payload.SourceMatcher,
		HighSeverity: severity.Severity(payload.HighSeverity), LowSeverity: severity.Severity(payload.LowSeverity), Enabled: payload.enabled(),
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, rule, requestID)
}

func (h *Handler) GetSuppression(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	rule, err := h.service.GetSuppression(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, rule, requestID)
}

func (h *Handler) ListSuppressions(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	pagination := httpx.ParsePagination(r.URL.Query())
	enabledOnly := r.URL.Query().Get("enabled") == "true"
	rules, total, err := h.service.ListSuppressions(r.Context(), enabledOnly, pagination.PerPage, pagination.Offset)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteMeta(w, http.StatusOK, rules, httpx.NewPageMeta(pagination, total), requestID)
}

func (h *Handler) DeleteSuppression(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	if err := h.service.DeleteSuppression(r.Context(), id); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusNoContent, nil, requestID)
}
