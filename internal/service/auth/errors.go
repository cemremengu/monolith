package auth

import "errors"

var (
	ErrSessionExpired       = errors.New("session expired")
	ErrSessionRevoked       = errors.New("session revoked")
	ErrSessionNeedsRotation = errors.New("session needs rotation")
	ErrInvalidSessionID     = errors.New("invalid session ID")
	ErrSessionNotFound      = errors.New("session not found")
)
