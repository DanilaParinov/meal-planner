package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateAPIKey returns a cryptographically random 64-character hex string
func GenerateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate API key: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// GenerateDeviceID returns a short random hex string suitable for identifying a user record
func GenerateDeviceID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate device ID: %w", err)
	}
	return "user-" + hex.EncodeToString(b), nil
}