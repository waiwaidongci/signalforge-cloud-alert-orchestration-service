package application

import (
	"context"
	"testing"
	"time"

	"github.com/acme/signalforge/internal/escalation/domain"
	"github.com/acme/signalforge/internal/shared/clock"
)

func TestCreateManyRejectsEmptyName(t *testing.T) {
	service := NewService(&fakeEscalationRepository{}, clock.FixedClock{Time: time.Unix(0, 0).UTC()})
	_, err := service.CreateMany(context.Background(), []domain.Policy{{Name: ""}})
	if err == nil {
		t.Fatal("expected empty name validation error")
	}
}
