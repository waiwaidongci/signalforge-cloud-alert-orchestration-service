package http

import (
	"log/slog"
	"net/http"

	notificationapp "github.com/acme/signalforge/internal/notification/application"
	notificationdomain "github.com/acme/signalforge/internal/notification/domain"
	"github.com/acme/signalforge/internal/shared/httpx"
)

type Handler struct {
	service *notificationapp.Service
	logger  *slog.Logger
}

func NewHandler(service *notificationapp.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	notification, err := h.service.Get(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, notification, requestID)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	pagination := httpx.ParsePagination(r.URL.Query())
	channel := r.URL.Query().Get("channel")
	status := notificationdomain.Status(r.URL.Query().Get("status"))
	notifications, total, err := h.service.List(r.Context(), channel, status, pagination.PerPage, pagination.Offset)
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteMeta(w, http.StatusOK, notifications, httpx.NewPageMeta(pagination, total), requestID)
}
