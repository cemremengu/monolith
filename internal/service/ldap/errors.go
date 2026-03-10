package ldap

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid LDAP credentials")
	ErrUserNotFound       = errors.New("LDAP user not found")
)
