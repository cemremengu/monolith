package account

import "errors"

var (
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("current password is incorrect")
)
