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

func TestValidateRefreshToken(t *testing.T) {
	ts := NewTokenService()

	// Test valid token
	token, err := ts.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("Expected no error generating token, got: %v", err)
	}

	validatedToken, err := ts.ValidateRefreshToken(token)
	if err != nil {
		t.Errorf("Expected no error validating token, got: %v", err)
	}

	if validatedToken != token {
		t.Error("Expected validated token to match original")
	}

	// Test invalid token (wrong length)
	_, err = ts.ValidateRefreshToken("invalid")
	if err == nil {
		t.Error("Expected error for invalid token")
	}

	// Test invalid token (not hex)
	_, err = ts.ValidateRefreshToken("gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg")
	if err == nil {
		t.Error("Expected error for non-hex token")
	}
}
