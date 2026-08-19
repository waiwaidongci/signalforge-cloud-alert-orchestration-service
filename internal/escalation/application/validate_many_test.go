package application

import (
	"testing"

	"github.com/acme/signalforge/internal/escalation/domain"
)

func TestValidateManyRejectsEmptyName(t *testing.T) {
	service := NewService(nil, nil)
	if err := service.ValidateMany([]domain.Policy{{Name: ""}}); err == nil {
		t.Fatal("expected empty name validation error")
	}
}
