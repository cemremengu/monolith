package session

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
	ID        uuid.UUID  `db:"id"`
	Token     string     `db:"token"`
	PrevToken *string    `db:"prev_token"`
	AccountID uuid.UUID  `db:"account_id"`
	UserAgent string     `db:"user_agent"`
	ClientIP  string     `db:"client_ip"`
	TokenSeen bool       `db:"token_seen"`
	SeenAt    *time.Time `db:"seen_at"`
	CreatedAt time.Time  `db:"created_at"`
	RotatedAt time.Time  `db:"rotated_at"`
	RevokedAt *time.Time `db:"revoked_at"`

	// UnhashedToken is used to store the unhashed token temporarily
	UnhashedToken string `db:"-" json:"-"`
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
