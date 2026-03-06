package order

import (
	"crypto/rand"
	"encoding/hex"
)

func generateOrderId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return "ORD_" + hex.EncodeToString(b)
}
