package http

import (
	"context"
	"log/slog"
	"net/http"

	incidentapp "github.com/acme/signalforge/internal/incident/application"
	incidentdomain "github.com/acme/signalforge/internal/incident/domain"
	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/httpx"
)

type Handler struct {
	service *incidentapp.Service
	logger  *slog.Logger
}

func NewHandler(service *incidentapp.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	pagination := httpx.ParsePagination(r.URL.Query())
	filter := incidentdomain.ListFilter{
		Fingerprint: r.URL.Query().Get("fingerprint"),
		Status:      incidentdomain.Status(r.URL.Query().Get("status")),
		Severity:    r.URL.Query().Get("severity"),
		Query:       r.URL.Query().Get("q"),
	}
	incidents, total, err := h.service.List(r.Context(), filter, pagination.PerPage, pagination.Offset)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteMeta(w, http.StatusOK, incidents, httpx.NewPageMeta(pagination, total), requestID)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	incident, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, incident, requestID)
}

func (h *Handler) Acknowledge(w http.ResponseWriter, r *http.Request) {
	h.batch(w, r, h.service.Acknowledge)
}

func (h *Handler) Close(w http.ResponseWriter, r *http.Request) {
	h.batch(w, r, h.service.Close)
}

func (h *Handler) Timeline(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	pagination := httpx.ParsePagination(r.URL.Query())
	events, total, err := h.service.Timeline(r.Context(), id, pagination.PerPage, pagination.Offset)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteMeta(w, http.StatusOK, detachTimeline(events), httpx.NewPageMeta(pagination, total), requestID)
}

func (h *Handler) batch(w http.ResponseWriter, r *http.Request, action func(ctx context.Context, ids []string, actor string) error) {
	requestID := httpx.RequestID(r.Context())
	var payload batchIncidentRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	if len(payload.IDs) == 0 {
		httpx.WriteError(w, requestID, apperr.BadRequest("IDS_REQUIRED", "ids 不能为空"))
		return
	}
	if err := action(r.Context(), payload.IDs, payload.Actor); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"affected": len(payload.IDs)}, requestID)
}
