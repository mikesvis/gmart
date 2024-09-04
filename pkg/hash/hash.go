package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
