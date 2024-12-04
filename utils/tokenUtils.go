package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
} 