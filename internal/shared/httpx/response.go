package httpx

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Data      any       `json:"data,omitempty"`
	Meta      any       `json:"meta,omitempty"`
	Error     *APIError `json:"error,omitempty"`
	RequestID string    `json:"request_id"`
}

type APIError struct {
	Code        string       `json:"code"`
	Message     string       `json:"message"`
	FieldErrors []FieldError `json:"field_errors,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, value any, requestID string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Data: value, RequestID: requestID})
}

func WriteMeta(w http.ResponseWriter, status int, data any, meta any, requestID string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Data: data, Meta: meta, RequestID: requestID})
}

func WriteRawJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
