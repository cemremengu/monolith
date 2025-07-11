package session

import (
	"crypto/sha256"
	"encoding/hex"

	"monolith/internal/util"
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

	return token, hashToken(token, secretKey), nil
}

func createToken() (string, error) {
	token, err := util.RandomHex(sessionTokenLength)
	if err != nil {
		return "", err
	}

	return token, nil
}
