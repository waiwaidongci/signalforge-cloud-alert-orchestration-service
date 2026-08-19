package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func New() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Sprintf("random id: %v", err))
	}
	return hex.EncodeToString(b[:])
}

func Prefix(prefix string) string {
	return prefix + "_" + New()
}
