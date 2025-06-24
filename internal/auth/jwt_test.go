package auth

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenService_GenerateAccessToken(t *testing.T) {
	ts := NewTokenService()

	userID := "test-user-id"
	email := "test@example.com"
	isAdmin := true

	token, err := ts.GenerateAccessToken(userID, email, isAdmin)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.IsType(t, "", token)
}

func TestTokenService_GenerateRefreshToken(t *testing.T) {
	ts := NewTokenService()

	userID := "test-user-id"

	token, err := ts.GenerateRefreshToken(userID)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.IsType(t, "", token)
}

func TestTokenService_ValidateToken_ValidToken(t *testing.T) {
	ts := NewTokenService()

	userID := "test-user-id"
	email := "test@example.com"
	isAdmin := true

	token, err := ts.GenerateAccessToken(userID, email, isAdmin)
	require.NoError(t, err)

	claims, err := ts.ValidateToken(token)

	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, isAdmin, claims.IsAdmin)
	assert.Equal(t, "monolith", claims.Issuer)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestTokenService_ValidateToken_InvalidToken(t *testing.T) {
	ts := NewTokenService()

	testCases := []struct {
		name          string
		token         string
		expectedError error
	}{
		{
			name:          "malformed token",
			token:         "invalid.token.string",
			expectedError: ErrTokenMalformed,
		},
		{
			name:          "empty token",
			token:         "",
			expectedError: ErrTokenMalformed,
		},
		{
			name:          "random string",
			token:         "random-string",
			expectedError: ErrTokenMalformed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := ts.ValidateToken(tc.token)

			assert.Error(t, err)
			assert.Nil(t, claims)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTokenService_ValidateRefreshToken_ValidToken(t *testing.T) {
	ts := NewTokenService()

	userID := "test-user-id"

	token, err := ts.GenerateRefreshToken(userID)
	require.NoError(t, err)

	claims, err := ts.ValidateRefreshToken(token)

	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "monolith", claims.Issuer)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestTokenService_ValidateRefreshToken_InvalidToken(t *testing.T) {
	ts := NewTokenService()

	testCases := []struct {
		name          string
		token         string
		expectedError error
	}{
		{
			name:          "malformed token",
			token:         "invalid.token.string",
			expectedError: ErrTokenMalformed,
		},
		{
			name:          "empty token",
			token:         "",
			expectedError: ErrTokenMalformed,
		},
		{
			name:          "random string",
			token:         "random-string",
			expectedError: ErrTokenMalformed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := ts.ValidateRefreshToken(tc.token)

			assert.Error(t, err)
			assert.Nil(t, claims)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTokenService_ExpiredToken(t *testing.T) {
	// Set very short token duration for testing
	os.Setenv("JWT_ACCESS_TOKEN_DURATION", "1ns")
	defer os.Unsetenv("JWT_ACCESS_TOKEN_DURATION")

	ts := NewTokenService()

	userID := "test-user-id"
	email := "test@example.com"
	isAdmin := false

	token, err := ts.GenerateAccessToken(userID, email, isAdmin)
	require.NoError(t, err)

	// Wait a moment to ensure token expires
	time.Sleep(time.Millisecond)

	claims, err := ts.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrTokenExpired, err)
}

func TestTokenService_ExpiredRefreshToken(t *testing.T) {
	// Set very short token duration for testing
	os.Setenv("JWT_REFRESH_TOKEN_DURATION", "1ns")
	defer os.Unsetenv("JWT_REFRESH_TOKEN_DURATION")

	ts := NewTokenService()

	userID := "test-user-id"

	token, err := ts.GenerateRefreshToken(userID)
	require.NoError(t, err)

	// Wait a moment to ensure token expires
	time.Sleep(time.Millisecond)

	claims, err := ts.ValidateRefreshToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrTokenExpired, err)
}

func TestTokenService_CrossValidation(t *testing.T) {
	ts := NewTokenService()

	userID := "test-user-id"
	email := "test@example.com"

	// Generate access token
	accessToken, err := ts.GenerateAccessToken(userID, email, false)
	require.NoError(t, err)

	// Generate refresh token
	refreshToken, err := ts.GenerateRefreshToken(userID)
	require.NoError(t, err)

	// Both tokens should be valid with their respective validation methods
	// Access token should be valid as access token
	accessClaims, err := ts.ValidateToken(accessToken)
	assert.NoError(t, err)
	assert.NotNil(t, accessClaims)
	assert.Equal(t, userID, accessClaims.UserID)

	// Refresh token should be valid as refresh token
	refreshClaims, err := ts.ValidateRefreshToken(refreshToken)
	assert.NoError(t, err)
	assert.NotNil(t, refreshClaims)
	assert.Equal(t, userID, refreshClaims.UserID)
}

func TestTokenService_DifferentSecrets(t *testing.T) {
	// Create first token service with default secret
	ts1 := NewTokenService()

	// Set different secret and create second token service
	os.Setenv("SECRET_KEY", "different-secret-key")
	defer os.Unsetenv("SECRET_KEY")
	ts2 := NewTokenService()

	userID := "test-user-id"
	email := "test@example.com"

	// Generate token with first service
	token, err := ts1.GenerateAccessToken(userID, email, false)
	require.NoError(t, err)

	// Try to validate with second service (should fail due to different secret)
	claims, err := ts2.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrTokenInvalid, err)
}

func TestHashToken(t *testing.T) {
	token := "test-token"
	hash1 := HashToken(token)
	hash2 := HashToken(token)

	assert.NotEmpty(t, hash1)
	assert.Equal(t, hash1, hash2) // Same input should produce same hash
	assert.Len(t, hash1, 64)      // SHA256 produces 64 character hex string

	// Different tokens should produce different hashes
	differentToken := "different-token"
	differentHash := HashToken(differentToken)
	assert.NotEqual(t, hash1, differentHash)
}

func TestTokenService_TokenDurations(t *testing.T) {
	ts := NewTokenService()

	// Test default durations
	assert.Equal(t, 10*time.Minute, ts.AccessTokenDuration())
	assert.Equal(t, 7*24*time.Hour, ts.RefreshTokenDuration())
}

func TestTokenService_CustomDurations(t *testing.T) {
	// Set custom durations
	os.Setenv("JWT_ACCESS_TOKEN_DURATION", "30m")
	os.Setenv("JWT_REFRESH_TOKEN_DURATION", "168h") // 7 days in hours
	defer func() {
		os.Unsetenv("JWT_ACCESS_TOKEN_DURATION")
		os.Unsetenv("JWT_REFRESH_TOKEN_DURATION")
	}()

	ts := NewTokenService()

	assert.Equal(t, 30*time.Minute, ts.AccessTokenDuration())
	assert.Equal(t, 168*time.Hour, ts.RefreshTokenDuration())
}