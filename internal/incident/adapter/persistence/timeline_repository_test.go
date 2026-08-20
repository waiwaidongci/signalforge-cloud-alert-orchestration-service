package persistence

import (
	"errors"
	"strings"
	"testing"

	"github.com/acme/signalforge/internal/shared/apperr"
)

var errIterate = errors.New("rows iteration failed")

type closeRows struct {
	iterateErr error
	closed     bool
}

func (r *closeRows) Next() bool      { return false }
func (*closeRows) Scan(...any) error { return nil }
func (r *closeRows) Err() error      { return r.iterateErr }
func (r *closeRows) Close() error    { r.closed = true; return errClose{} }

type errClose struct{}

func (errClose) Error() string { return "rows close failed" }
func TestCollectTimelineRowsPropagatesCloseError(t *testing.T) {
	rows := &closeRows{}
	_, err := (&TimelineSQLRepository{}).collectRows(rows)
	if err == nil || !strings.Contains(err.Error(), "rows close failed") || !strings.Contains(err.Error(), "timeline repository") {
		t.Fatalf("err=%v", err)
	}
	if !rows.closed {
		t.Fatal("rows were not closed")
	}
}

func TestCollectTimelineRowsClosesAfterIterationError(t *testing.T) {
	rows := &closeRows{iterateErr: errIterate}
	_, err := (&TimelineSQLRepository{}).collectRows(rows)
	if !errors.Is(err, errIterate) || !strings.Contains(err.Error(), "rows close failed") || !strings.Contains(err.Error(), "timeline repository") || !apperr.IsKind(err, apperr.KindInternal) {
		t.Fatalf("err=%v", err)
	}
}

func TestMergeTimelineCloseErrorKeepsPrimaryWhenCloseSucceeds(t *testing.T) {
	primary := errors.New("scan failed")
	if !errors.Is(mergeTimelineCloseError(primary, nil), primary) {
		t.Fatal("primary error was lost")
	}
}
