package auth

import (
	"time"

	"monolith/internal/service/user"
)

type LoginRequest struct {
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  user.Account `json:"user"`
}

type SessionResponse struct {
	SessionID  string    `json:"sessionId"`
	DeviceInfo string    `json:"deviceInfo"`
	IPAddress  string    `json:"ipAddress"`
	CreatedAt  time.Time `json:"createdAt"`
	RotatedAt  time.Time `json:"rotatedAt"`
	IsCurrent  bool      `json:"isCurrent"`
}
