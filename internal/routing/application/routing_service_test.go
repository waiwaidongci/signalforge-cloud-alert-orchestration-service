package application

import (
	"context"
	"testing"

	"github.com/acme/signalforge/internal/routing/domain"
	"github.com/acme/signalforge/internal/shared/matcher"
)

type fakeRoutingRepository struct {
	rules []domain.Rule
}

func (f *fakeRoutingRepository) Active(context.Context) ([]domain.Rule, error) {
	return f.rules, nil
}

func (f *fakeRoutingRepository) Create(context.Context, domain.Rule) error { return nil }
func (f *fakeRoutingRepository) Update(context.Context, domain.Rule) error { return nil }
func (f *fakeRoutingRepository) FindByID(context.Context, string) (domain.Rule, error) {
	return domain.Rule{}, nil
}
func (f *fakeRoutingRepository) List(context.Context, bool, int, int) ([]domain.Rule, int, error) {
	return f.rules, 0, nil
}
func (f *fakeRoutingRepository) Delete(context.Context, string) error { return nil }

func TestResolvePreservesRepositoryState(t *testing.T) {
	repo := &fakeRoutingRepository{rules: []domain.Rule{
		{Name: "low", Priority: 1, Enabled: true, Match: matcher.Selector{}, Channels: []domain.Channel{{Channel: "email", Destination: "low@example.com"}}},
		{Name: "high", Priority: 10, Enabled: true, Match: matcher.Selector{}, Channels: []domain.Channel{{Channel: "webhook", Destination: "https://example.com"}}},
	}}
	service := NewService(repo, nil)
	first, err := service.Resolve(context.Background(), matcher.Target{})
	if err != nil || len(first) != 1 || first[0].Destination != "https://example.com" {
		t.Fatalf("unexpected first result: %+v, %v", first, err)
	}
	first[0].Destination = "mutated"
	second, err := service.Resolve(context.Background(), matcher.Target{})
	if err != nil || len(second) != 1 || second[0].Destination != "https://example.com" {
		t.Fatalf("repository state was mutated: first=%+v second=%+v err=%v", first, second, err)
	}
	if repo.rules[1].Channels[0].Destination != "https://example.com" {
		t.Fatalf("repository channels were mutated: %+v", repo.rules[1].Channels)
	}
}
