// Cryptographic utilities for password management
package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"  
)

// HashPassword generates a bcrypt hash from a plaintext password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
} 

// CheckPassword compares a plaintext password with a hashed password
func CheckPassword(password string, hashedPassword string) error {
	// Returns an error if the passwords do not match
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
