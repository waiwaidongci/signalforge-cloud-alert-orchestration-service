package http

import (
	"log/slog"
	"net/http"

	"github.com/acme/signalforge/internal/shared/httpx"
	sourceapp "github.com/acme/signalforge/internal/source/application"
	sourcedomain "github.com/acme/signalforge/internal/source/domain"
)

type Handler struct {
	service *sourceapp.Service
	logger  *slog.Logger
}

func NewHandler(service *sourceapp.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	var payload sourceRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	source, err := h.service.Create(r.Context(), sourcedomain.Source{
		Name: payload.Name, Kind: sourcedomain.Kind(payload.Kind), Enabled: payload.enabled(), RateLimit: payload.RateLimit, FieldMapping: payload.FieldMapping,
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, sourceResponse{
		ID: source.ID, Name: source.Name, Kind: string(source.Kind), APIKey: source.APIKey, Enabled: source.Enabled, RateLimit: source.RateLimit, FieldMapping: source.FieldMapping, CreatedAt: source.CreatedAt, UpdatedAt: source.UpdatedAt,
	}, requestID)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	var payload sourceRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	source, err := h.service.Update(r.Context(), sourcedomain.Source{
		ID: id, Name: payload.Name, Kind: sourcedomain.Kind(payload.Kind), Enabled: payload.enabled(), RateLimit: payload.RateLimit, FieldMapping: payload.FieldMapping,
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, sourceResponse{
		ID: source.ID, Name: source.Name, Kind: string(source.Kind), APIKey: source.APIKey, Enabled: source.Enabled, RateLimit: source.RateLimit, FieldMapping: source.FieldMapping, CreatedAt: source.CreatedAt, UpdatedAt: source.UpdatedAt,
	}, requestID)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	source, err := h.service.Get(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, sourceResponse{
		ID: source.ID, Name: source.Name, Kind: string(source.Kind), APIKey: source.APIKey, Enabled: source.Enabled, RateLimit: source.RateLimit, FieldMapping: source.FieldMapping, CreatedAt: source.CreatedAt, UpdatedAt: source.UpdatedAt,
	}, requestID)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	pagination := httpx.ParsePagination(r.URL.Query())
	var enabled *bool
	if raw := r.URL.Query().Get("enabled"); raw != "" {
		value := raw == "true" || raw == "1"
		enabled = &value
	}
	sources, total, err := h.service.List(r.Context(), enabled, pagination.PerPage, pagination.Offset)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	response := make([]sourceResponse, 0, len(sources))
	for _, source := range sources {
		response = append(response, sourceResponse{
			ID: source.ID, Name: source.Name, Kind: string(source.Kind), APIKey: source.APIKey, Enabled: source.Enabled, RateLimit: source.RateLimit, FieldMapping: source.FieldMapping, CreatedAt: source.CreatedAt, UpdatedAt: source.UpdatedAt,
		})
	}
	httpx.WriteMeta(w, http.StatusOK, response, httpx.NewPageMeta(pagination, total), requestID)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusNoContent, nil, requestID)
}
