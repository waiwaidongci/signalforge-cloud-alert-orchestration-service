package http

import (
	"errors"
	"log/slog"
	"net/http"

	escalationapp "github.com/acme/signalforge/internal/escalation/application"
	escalationdomain "github.com/acme/signalforge/internal/escalation/domain"
	"github.com/acme/signalforge/internal/shared/apperr"
	"github.com/acme/signalforge/internal/shared/httpx"
)

type Handler struct {
	service *escalationapp.Service
	logger  *slog.Logger
}

func NewHandler(service *escalationapp.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	var payload escalationRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	policy, err := h.service.Create(r.Context(), escalationdomain.Policy{
		Name: payload.Name, Description: payload.Description, Matcher: payload.Matcher, WaitSeconds: payload.WaitSeconds,
		RepeatSeconds: payload.RepeatSeconds, MaxRepeats: payload.MaxRepeats, Routes: toRoutes(payload.Routes), Enabled: payload.enabled(),
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, policy, requestID)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	var payload escalationRequest
	if err := httpx.DecodeJSON(r, &payload); err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	policy, err := h.service.Update(r.Context(), escalationdomain.Policy{
		ID: id, Name: payload.Name, Description: payload.Description, Matcher: payload.Matcher, WaitSeconds: payload.WaitSeconds,
		RepeatSeconds: payload.RepeatSeconds, MaxRepeats: payload.MaxRepeats, Routes: toRoutes(payload.Routes), Enabled: payload.enabled(),
	})
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, policy, requestID)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	id, err := httpx.RequirePathValue(r, "id")
	if err != nil {
		httpx.WriteError(w, requestID, err)
		return
	}
	policy, err := h.service.Get(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, requestID, policyReadError(err))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, policy, requestID)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	requestID := httpx.RequestID(r.Context())
	pagination := httpx.ParsePagination(r.URL.Query())
	enabledOnly := r.URL.Query().Get("enabled") == "true"
	policies, total, err := h.service.List(r.Context(), enabledOnly, pagination.PerPage, pagination.Offset)
	if err != nil {
		httpx.WriteError(w, requestID, policyReadError(err))
		return
	}
	httpx.WriteMeta(w, http.StatusOK, policies, httpx.NewPageMeta(pagination, total), requestID)
}

func policyReadError(err error) error {
	if errors.Is(err, escalationdomain.ErrInvalidPolicyData) {
		return apperr.Wrap(err, apperr.KindBadRequest, "MALFORMED_ESCALATION_POLICY", "升级策略数据无效")
	}
	return err
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

func toRoutes(values []routeRequest) []escalationdomain.Route {
	routes := make([]escalationdomain.Route, 0, len(values))
	for _, value := range values {
		routes = append(routes, escalationdomain.Route{Channel: value.Channel, Destination: value.Destination, Severities: value.Severities})
	}
	return routes
}
