package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"monolith/internal/config"
)

const sessionTokenLength = 32

type TokenService struct {
	config *config.SecurityConfig
}

func NewTokenService() *TokenService {
	return &TokenService{
		config: config.NewSecurityConfig(),
	}
}

func hashToken(token string, secretKey string) string {
	hash := sha256.Sum256([]byte(token + secretKey))
	return hex.EncodeToString(hash[:])
}

func (ts *TokenService) GenerateSessionToken() (string, error) {
	bytes := make([]byte, sessionTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
