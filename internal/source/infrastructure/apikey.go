package infrastructure

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func NewAPIKey() string {
	var b [24]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Sprintf("generate api key: %v", err))
	}
	return "sf_" + hex.EncodeToString(b[:])
}
