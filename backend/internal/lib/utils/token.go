package utils

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	DEFAULT_TOKEN_LENGTH = 12
)

func GenerateToken() string {
	bytes := make([]byte, DEFAULT_TOKEN_LENGTH)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
