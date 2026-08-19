package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/acme/signalforge/internal/escalation/domain"
	"github.com/acme/signalforge/internal/shared/clock"
)

type fakeEscalationRepository struct {
	created []domain.Policy
}

func (f *fakeEscalationRepository) Create(_ context.Context, p domain.Policy) error {
	if p.Name == "fail" {
		return errors.New("create failed")
	}
	f.created = append(f.created, p)
	return nil
}
func (f *fakeEscalationRepository) Update(context.Context, domain.Policy) error { return nil }
func (f *fakeEscalationRepository) FindByID(context.Context, string) (domain.Policy, error) {
	return domain.Policy{}, nil
}
func (f *fakeEscalationRepository) List(context.Context, bool, int, int) ([]domain.Policy, int, error) {
	return f.created, len(f.created), nil
}
func (f *fakeEscalationRepository) Delete(context.Context, string) error { return nil }
func (f *fakeEscalationRepository) Active(context.Context) ([]domain.Policy, error) {
	return f.created, nil
}

func TestCreateManyReturnsErrorOnFailure(t *testing.T) {
	service := NewService(&fakeEscalationRepository{}, clock.FixedClock{Time: time.Unix(0, 0).UTC()})
	_, err := service.CreateMany(context.Background(), []domain.Policy{
		{Name: "ok"},
		{Name: "fail"},
	})
	if err == nil {
		t.Fatal("expected batch error")
	}
}
