package matcher

import (
	"errors"
	"testing"
)

func TestDecodeSelectorRejectsNull(t *testing.T) {
	_, err := DecodeSelector("null")
	if !errors.Is(err, ErrInvalidSelector) {
		t.Fatalf("err=%v", err)
	}
}
