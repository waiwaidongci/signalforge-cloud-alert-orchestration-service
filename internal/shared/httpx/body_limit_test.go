package httpx

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/acme/signalforge/internal/shared/apperr"
)

func TestDecodeJSONPreservesRequestTooLargeError(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/alerts", nil)
	request.Body = http.MaxBytesReader(recorder, io.NopCloser(strings.NewReader(`{"message":"too large"}`)), 8)

	var payload map[string]string
	err := DecodeJSON(request, &payload)
	var tooLarge *http.MaxBytesError
	if !errors.As(err, &tooLarge) {
		t.Fatalf("expected MaxBytesError in chain, got %v", err)
	}
}

func TestWriteErrorMapsRequestTooLargeTo413(t *testing.T) {
	recorder := httptest.NewRecorder()
	WriteError(recorder, "request-17", apperr.PayloadTooLarge(errors.New("body limit")))
	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("got status %d, want %d", recorder.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestMaxBodyBytesReturnsPayloadTooLargeResponse(t *testing.T) {
	handler := MaxBodyBytes(8, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]string
		if err := DecodeJSON(r, &payload); err != nil {
			WriteError(w, "request-17", err)
			return
		}
		WriteJSON(w, http.StatusCreated, payload, "request-17")
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/alerts", strings.NewReader(`{"message":"too large"}`)))
	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("got status %d, want %d", recorder.Code, http.StatusRequestEntityTooLarge)
	}
}
