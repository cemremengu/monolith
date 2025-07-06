package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"monolith/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenInvalid     = errors.New("token invalid")
	ErrTokenMalformed   = errors.New("token malformed")
	ErrTokenNotValidYet = errors.New("token not valid yet")
)

type Claims struct {
	UserID  string `json:"userId"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

const refreshTokenLength = 32

type TokenService struct {
	config *config.SecurityConfig
}

func NewTokenService() *TokenService {
	return &TokenService{
		config: config.NewSecurityConfig(),
	}
}

func (ts *TokenService) getJWTSecret() []byte {
	return []byte(ts.config.SecretKey)
}

func HashToken(token string, secretKey string) string {
	hash := sha256.Sum256([]byte(token + secretKey))
	return hex.EncodeToString(hash[:])
}

func (ts *TokenService) GenerateAccessToken(userID, email string, isAdmin bool) (string, error) {
	claims := Claims{
		UserID:  userID,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ts.config.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "monolith",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.getJWTSecret())
}

func (ts *TokenService) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, refreshTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (ts *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return ts.getJWTSecret(), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrTokenNotValidYet
		default:
			return nil, ErrTokenInvalid
		}
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

func (ts *TokenService) AccessTokenDuration() time.Duration {
	return ts.config.AccessTokenDuration
}

func (ts *TokenService) RefreshTokenDuration() time.Duration {
	return ts.config.RefreshTokenDuration
}
