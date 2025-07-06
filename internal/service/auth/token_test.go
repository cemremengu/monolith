package auth

import (
	"testing"
)

func TestGenerateRefreshToken(t *testing.T) {
	ts := NewTokenService()

	// Test basic generation
	token, err := ts.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check token length (32 bytes = 64 hex characters)
	expectedLength := refreshTokenLength * 2
	if len(token) != expectedLength {
		t.Errorf("Expected token length %d, got %d", expectedLength, len(token))
	}

	// Test multiple tokens are different
	token2, err := ts.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if token == token2 {
		t.Error("Expected different tokens, got identical tokens")
	}
}
