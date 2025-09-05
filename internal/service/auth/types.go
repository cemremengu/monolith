package auth

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	ID        uuid.UUID `json:"id"`
	UserAgent string    `json:"userAgent"`
	ClientIP  string    `json:"clientIp"`
	CreatedAt time.Time `json:"createdAt"`
	RotatedAt time.Time `json:"rotatedAt"`
	IsCurrent bool      `json:"isCurrent"`
}

type Session struct {
	ID        uuid.UUID
	Token     string
	PrevToken *string
	AccountID uuid.UUID
	UserAgent string
	ClientIP  string
	TokenSeen bool
	SeenAt    *time.Time
	CreatedAt time.Time
	RotatedAt time.Time
	RevokedAt *time.Time

	// UnhashedToken is used to store the unhashed token temporarily
	UnhashedToken string `json:"-"`
}

const (
	rotationLeeway   = 5 * time.Second
	urgentRotateTime = 1 * time.Minute
)

func (s *Session) NextRotation(rotationInterval time.Duration) time.Time {
	return s.RotatedAt.Add(rotationInterval - rotationLeeway)
}

func (s *Session) NeedsRotation(rotationInterval time.Duration) bool {
	if !s.TokenSeen {
		return s.RotatedAt.Before(time.Now().Add(-urgentRotateTime))
	}

	return s.RotatedAt.Before(time.Now().Add(-rotationInterval))
}

type CreateSessionRequest struct {
	AccountID uuid.UUID
	ClientIP  string
	UserAgent string
}

type RotateSessionRequest struct {
	UnhashedToken string
	ClientIP      string
	UserAgent     string
}

// AuthUser holds the authenticated user information stored in the context
type AuthUser struct {
	UserID      uuid.UUID
	UserEmail   string
	IsAdmin     bool
	SessionID   uuid.UUID
	WorkspaceID uuid.UUID
}
