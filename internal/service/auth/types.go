package auth

import (
	"time"

	"github.com/google/uuid"
)

// Account status constants
const (
	// AccountStatusActive indicates the account is fully active and can authenticate
	AccountStatusActive = "active"
	// AccountStatusPending indicates the account is created but not yet activated (e.g., awaiting email verification)
	AccountStatusPending = "pending"
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
	AccountID   uuid.UUID
	Email       string
	IsAdmin     bool
	SessionID   uuid.UUID
	WorkspaceID uuid.UUID
}

// AuthContext holds the complete authentication context from a consolidated query
type AuthContext struct {
	SessionID      uuid.UUID
	SessionToken   string
	AccountID      uuid.UUID
	AccountEmail   string
	AccountIsAdmin bool
	AccountStatus  string
	WorkspaceID    uuid.UUID
	SessionCreated time.Time
	SessionRotated time.Time
	SessionRevoked *time.Time
}
