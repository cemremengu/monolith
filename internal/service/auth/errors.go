package auth

import "errors"

var (
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrSessionExpired      = errors.New("session expired or invalid")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrInvalidSessionID    = errors.New("invalid session ID")
	ErrSessionNotFound     = errors.New("session not found")
)
