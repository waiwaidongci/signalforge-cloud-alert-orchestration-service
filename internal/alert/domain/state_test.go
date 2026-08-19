package domain

import "testing"

func TestAcknowledgedStatusIsValid(t *testing.T) {
	if !StatusAcknowledged.Valid() {
		t.Fatal("acknowledged must be a valid alert status")
	}
}
