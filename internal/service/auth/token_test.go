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

	// Test multiple tokens are different
	token2, err := ts.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if token == token2 {
		t.Error("Expected different tokens, got identical tokens")
	}
}
