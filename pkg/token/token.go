package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// GenerateInvitationToken generates a secure 32-character token for invitations
func GenerateInvitationToken() (string, error) {
	// Generate 24 bytes (32 base64 characters)
	bytes := make([]byte, 24)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	// Use URL-safe base64 encoding and take first 32 characters
	token := base64.URLEncoding.EncodeToString(bytes)
	// Remove padding if any and ensure length
	if len(token) > 32 {
		token = token[:32]
	}
	return token, nil
}

// GenerateOTP generates a 6-digit numeric OTP code (000000-999999)
func GenerateOTP() (string, error) {
	// Generate random number between 0 and 999999
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}
	// Format as 6-digit string with leading zeros
	return fmt.Sprintf("%06d", n.Int64()), nil
}

