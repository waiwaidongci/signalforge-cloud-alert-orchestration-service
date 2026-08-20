package store

import "testing"

func TestDecodePayloadInitializesEmptyMap(t *testing.T) {
	payload, err := DecodePayload("null")
	if payload != nil || !IsEmptyPayload(err) {
		t.Fatal("null payload must be classified explicitly")
	}
}
