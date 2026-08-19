package infrastructure

import (
	"crypto/sha256"
	"encoding/hex"
)

func ShortHash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:8])
}
