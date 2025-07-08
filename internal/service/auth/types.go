package auth

import (
	"time"

	"monolith/internal/service/account"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Token   string          `json:"token"`
	Account account.Account `json:"account"`
}

type SessionResponse struct {
	ID        uuid.UUID `json:"id"`
	UserAgent string    `json:"userAgent"`
	ClientIP  string    `json:"clientIp"`
	CreatedAt time.Time `json:"createdAt"`
	RotatedAt time.Time `json:"rotatedAt"`
	IsCurrent bool      `json:"isCurrent"`
}
