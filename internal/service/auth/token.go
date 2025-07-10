package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

const sessionTokenLength = 16

func hashToken(token string, secretKey string) string {
	hash := sha256.Sum256([]byte(token + secretKey))
	return hex.EncodeToString(hash[:])
}

func createAndHashToken(secretKey string) (string, string, error) {
	token, err := createToken()
	if err != nil {
		return "", "", err
	}

	return token, hashToken(secretKey, token), nil
}

func createToken() (string, error) {
	bytes := make([]byte, sessionTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
