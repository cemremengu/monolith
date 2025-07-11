package session

import "errors"

var (
	ErrSessionExpired   = errors.New("session expired")
	ErrSessionRevoked   = errors.New("session revoked")
	ErrInvalidSessionID = errors.New("invalid session ID")
	ErrSessionNotFound  = errors.New("session not found")
)
