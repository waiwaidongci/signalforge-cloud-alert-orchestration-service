package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/acme/signalforge/internal/shared/apperr"
)

func WriteError(w http.ResponseWriter, requestID string, err error) {
	var target *apperr.Error
	if !errors.As(err, &target) {
		slog.Error("http handler returned non-app error", "error", err)
		target = apperr.Internal("INTERNAL_ERROR", "服务器内部错误")
	} else if target.Kind == apperr.KindInternal {
		slog.Error("http handler returned internal error", "code", target.Code, "message", target.Message, "cause", target.Unwrap())
	}
	status := target.HTTPStatus
	if status == 0 {
		status = http.StatusInternalServerError
	}
	apiErr := APIError{Code: target.Code, Message: target.Message}
	for _, detail := range target.Details {
		apiErr.FieldErrors = append(apiErr.FieldErrors, FieldError{Field: detail.Field, Message: detail.Message})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Error: &apiErr, RequestID: requestID})
}
