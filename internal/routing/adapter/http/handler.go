package http

import (
	"log/slog"
	"net/http"

	routingapp "github.com/acme/signalforge/internal/routing/application"
	routingdomain "github.com/acme/signalforge/internal/routing/domain"
	"github.com/acme/signalforge/internal/shared/httpx"
)

type Handler struct {
	service *routingapp.Service
	logger  *slog.Logger
}

func NewHandler(service *routingapp.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	var payload routingRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	rule, err := h.service.Create(r.Context(), routingdomain.Rule{
		Name: payload.Name, Priority: payload.Priority, Match: payload.Match, Channels: toChannels(payload.Channels), Enabled: payload.enabled(),
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, rule, requestID)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	var payload routingRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	rule, err := h.service.Update(r.Context(), routingdomain.Rule{
		ID: id, Name: payload.Name, Priority: payload.Priority, Match: payload.Match, Channels: toChannels(payload.Channels), Enabled: payload.enabled(),
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, rule, requestID)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	rule, err := h.service.Get(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, rule, requestID)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	pagination := httpx.ParsePagination(r.URL.Query())
	enabledOnly := r.URL.Query().Get("enabled") == "true"
	rules, total, err := h.service.List(r.Context(), enabledOnly, pagination.PerPage, pagination.Offset)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteMeta(w, http.StatusOK, rules, httpx.NewPageMeta(pagination, total), requestID)
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

func toChannels(values []channelRequest) []routingdomain.Channel {
	channels := make([]routingdomain.Channel, 0, len(values))
	for _, value := range values {
		channels = append(channels, routingdomain.Channel{Channel: value.Channel, Destination: value.Destination})
	}
	return channels
}
