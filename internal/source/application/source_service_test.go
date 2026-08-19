package application

import (
	"context"
	"reflect"
	"testing"

	"github.com/acme/signalforge/internal/source/domain"
)

type fakeSourceRepositoryForFilter struct {
	sources []domain.Source
}

func (f *fakeSourceRepositoryForFilter) Create(context.Context, domain.Source) error { return nil }
func (f *fakeSourceRepositoryForFilter) Update(context.Context, domain.Source) error { return nil }
func (f *fakeSourceRepositoryForFilter) FindByID(context.Context, string) (domain.Source, error) { return domain.Source{}, nil }
func (f *fakeSourceRepositoryForFilter) FindByAPIKey(context.Context, string) (domain.Source, error) { return domain.Source{}, nil }
func (f *fakeSourceRepositoryForFilter) List(context.Context, *bool, int, int) ([]domain.Source, int, error) { return f.sources, len(f.sources), nil }
func (f *fakeSourceRepositoryForFilter) Delete(context.Context, string) error { return nil }

func TestListFilterDoesNotMutateRepository(t *testing.T) {
	enabled := true
	repo := &fakeSourceRepositoryForFilter{sources: []domain.Source{
		{Name: "one", Enabled: true},
		{Name: "two", Enabled: false},
		{Name: "three", Enabled: true},
	}}
	service := NewService(repo, nil)
	got, _, err := service.List(context.Background(), &enabled, 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 enabled sources, got %d", len(got))
	}
	if len(repo.sources) != 3 {
		t.Fatalf("repository slice was mutated, len=%d", len(repo.sources))
	}
	if !reflect.DeepEqual(repo.sources, []domain.Source{{Name: "one", Enabled: true}, {Name: "two", Enabled: false}, {Name: "three", Enabled: true}}) {
		t.Fatalf("repository contents were mutated: %+v", repo.sources)
	}
}
