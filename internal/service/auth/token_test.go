package auth

import (
	"testing"
)

func TestGenerateSessionToken(t *testing.T) {
	ts := NewTokenService()

	// Test basic generation
	token, err := ts.GenerateSessionToken()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Test multiple tokens are different
	token2, err := ts.GenerateSessionToken()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if token == token2 {
		t.Error("Expected different tokens, got identical tokens")
	}

	// Test token length (32 bytes = 64 hex characters)
	if len(token) != 64 {
		t.Errorf("Expected token length of 64 characters, got: %d", len(token))
	}
}
