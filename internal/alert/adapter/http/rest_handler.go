package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	alertapp "github.com/acme/signalforge/internal/alert/application"
	alertdomain "github.com/acme/signalforge/internal/alert/domain"
	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/httpx"
)

type RestHandler struct {
	service *alertapp.Service
	logger  *slog.Logger
}

func NewRestHandler(service *alertapp.Service, logger *slog.Logger) *RestHandler {
	return &RestHandler{service: service, logger: logger}
}

func (h *RestHandler) List(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	pagination := httpx.ParsePagination(r.URL.Query())
	filter := alertdomain.ListFilter{
		SourceID:    r.URL.Query().Get("source_id"),
		Fingerprint: r.URL.Query().Get("fingerprint"),
		Status:      alertdomain.Status(r.URL.Query().Get("status")),
		Severity:    r.URL.Query().Get("severity"),
		IncidentID:  r.URL.Query().Get("incident_id"),
		Query:       r.URL.Query().Get("q"),
	}
	alerts, total, err := h.service.List(r.Context(), filter, pagination.PerPage, pagination.Offset)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteMeta(w, http.StatusOK, alerts, httpx.NewPageMeta(pagination, total), requestID)
}

func (h *RestHandler) Get(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	alert, err := h.service.Get(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, alert, requestID)
}

func (h *RestHandler) BatchAcknowledge(w http.ResponseWriter, r *http.Request) {
	h.batch(w, r, h.service.Acknowledge)
}

func (h *RestHandler) BatchClose(w http.ResponseWriter, r *http.Request) {
	h.batch(w, r, h.service.Close)
}

func (h *RestHandler) batch(w http.ResponseWriter, r *http.Request, action func(ctx context.Context, ids []string, actor string) error) {
	requestID := httpx.RequestID(r.Context())
	var payload batchIDsRequest
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
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"acknowledged": len(payload.IDs)}, requestID)
}

func NormalizeQuery(raw string) string {
	return strings.TrimSpace(raw)
}
