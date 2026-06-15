package custom_bcrypt

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// hash pass
func HashPassword(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(passwordBytes), nil
}

// verify password
func VerifyPassword(password, passwordHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordHash)) == nil
}
