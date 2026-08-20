package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	escalationapp "github.com/acme/signalforge/internal/escalation/application"
	"github.com/acme/signalforge/internal/escalation/domain"
)

type malformedPolicyRepository struct{}

func (malformedPolicyRepository) Create(context.Context, domain.Policy) error { return nil }
func (malformedPolicyRepository) Update(context.Context, domain.Policy) error { return nil }
func (malformedPolicyRepository) FindByID(context.Context, string) (domain.Policy, error) {
	return domain.Policy{}, domain.ErrInvalidPolicyData
}
func (malformedPolicyRepository) List(context.Context, bool, int, int) ([]domain.Policy, int, error) {
	return nil, 0, nil
}
func (malformedPolicyRepository) Delete(context.Context, string) error            { return nil }
func (malformedPolicyRepository) Active(context.Context) ([]domain.Policy, error) { return nil, nil }

func TestGetMapsMalformedPolicyToBadRequest(t *testing.T) {
	handler := NewHandler(escalationapp.NewService(malformedPolicyRepository{}, nil), nil)
	request := httptest.NewRequest(http.MethodGet, "/policies/broken", nil)
	request.SetPathValue("id", "broken")
	response := httptest.NewRecorder()
	handler.Get(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
}
