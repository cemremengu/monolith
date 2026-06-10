package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashToken(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		secretKey string
	}{
		{
			name:      "basic token",
			token:     "mytoken123",
			secretKey: "mysecretkey",
		},
		{
			name:      "empty token",
			token:     "",
			secretKey: "mysecretkey",
		},
		{
			name:      "empty secret",
			token:     "mytoken123",
			secretKey: "",
		},
		{
			name:      "long token",
			token:     "abcdefghijklmnopqrstuvwxyz1234567890",
			secretKey: "secretkey123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1 := hashToken(tt.token, tt.secretKey)
			hash2 := hashToken(tt.token, tt.secretKey)

			assert.NotEmpty(t, hash1)
			assert.Equal(t, hash1, hash2, "hashing same token should produce same result")
			assert.Len(t, hash1, 64, "SHA256 hash should be 64 hex characters")
		})
	}
}

func TestHashToken_DifferentInputs(t *testing.T) {
	tests := []struct {
		name       string
		token1     string
		secretKey1 string
		token2     string
		secretKey2 string
	}{
		{
			name:       "different tokens same secret",
			token1:     "token1",
			secretKey1: "secret",
			token2:     "token2",
			secretKey2: "secret",
		},
		{
			name:       "same token different secrets",
			token1:     "token",
			secretKey1: "secret1",
			token2:     "token",
			secretKey2: "secret2",
		},
		{
			name:       "different tokens different secrets",
			token1:     "token1",
			secretKey1: "secret1",
			token2:     "token2",
			secretKey2: "secret2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1 := hashToken(tt.token1, tt.secretKey1)
			hash2 := hashToken(tt.token2, tt.secretKey2)

			assert.NotEqual(t, hash1, hash2, "different inputs should produce different hashes")
		})
	}
}

func TestCreateToken(t *testing.T) {
	t.Run("creates valid token", func(t *testing.T) {
		token, err := createToken()
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Len(t, token, sessionTokenLength*2, "token should be hex encoded (2 chars per byte)")
	})

	t.Run("creates unique tokens", func(t *testing.T) {
		token1, err1 := createToken()
		token2, err2 := createToken()

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2, "tokens should be unique")
	})
}

func TestCreateAndHashToken(t *testing.T) {
	secretKey := "test-secret-key"

	t.Run("returns token and hash", func(t *testing.T) {
		token, hashedToken, err := createAndHashToken(secretKey)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, hashedToken)
		assert.NotEqual(t, token, hashedToken)
	})

	t.Run("hash is deterministic for token", func(t *testing.T) {
		token, hashedToken, err := createAndHashToken(secretKey)
		require.NoError(t, err)

		rehashedToken := hashToken(token, secretKey)
		assert.Equal(t, hashedToken, rehashedToken)
	})

	t.Run("different calls produce different tokens", func(t *testing.T) {
		token1, hash1, err1 := createAndHashToken(secretKey)
		token2, hash2, err2 := createAndHashToken(secretKey)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
		assert.NotEqual(t, hash1, hash2)
	})
}

func TestRandomHex(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "length 8",
			length: 8,
		},
		{
			name:   "length 16",
			length: 16,
		},
		{
			name:   "length 32",
			length: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hex, err := randomHex(tt.length)
			require.NoError(t, err)
			assert.Len(t, hex, tt.length*2, "hex string should be double the byte length")
		})
	}

	t.Run("produces unique values", func(t *testing.T) {
		hex1, _ := randomHex(16)
		hex2, _ := randomHex(16)
		assert.NotEqual(t, hex1, hex2)
	})
}
