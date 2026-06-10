package login

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidUserID       = errors.New("invalid user ID")
)
