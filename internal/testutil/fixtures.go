package testutil

import (
	"time"

	"monolith/internal/service/account"
	"monolith/internal/service/auth"

	"github.com/google/uuid"
)

func NewTestAccount(opts ...func(*account.Account)) *account.Account {
	acc := &account.Account{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		Name:      StringPtr("Test User"),
		Password:  "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/X.xnWfOJZqjXX6K2K", // "password123"
		IsAdmin:   false,
		Language:  StringPtr("en"),
		Theme:     StringPtr("light"),
		Timezone:  StringPtr("UTC"),
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(acc)
	}

	return acc
}

func WithAccountID(id uuid.UUID) func(*account.Account) {
	return func(a *account.Account) {
		a.ID = id
	}
}

func WithUsername(username string) func(*account.Account) {
	return func(a *account.Account) {
		a.Username = username
	}
}

func WithEmail(email string) func(*account.Account) {
	return func(a *account.Account) {
		a.Email = email
	}
}

func WithPassword(password string) func(*account.Account) {
	return func(a *account.Account) {
		a.Password = password
	}
}

func WithIsAdmin(isAdmin bool) func(*account.Account) {
	return func(a *account.Account) {
		a.IsAdmin = isAdmin
	}
}

func WithStatus(status string) func(*account.Account) {
	return func(a *account.Account) {
		a.Status = status
	}
}

func NewTestSession(opts ...func(*auth.Session)) *auth.Session {
	now := time.Now()
	session := &auth.Session{
		ID:            uuid.New(),
		Token:         "hashed_token_123",
		AccountID:     uuid.New(),
		UserAgent:     "Mozilla/5.0 Test Browser",
		ClientIP:      "127.0.0.1",
		TokenSeen:     true,
		CreatedAt:     now,
		RotatedAt:     now,
		UnhashedToken: "unhashed_token_123",
	}

	for _, opt := range opts {
		opt(session)
	}

	return session
}

func WithSessionID(id uuid.UUID) func(*auth.Session) {
	return func(s *auth.Session) {
		s.ID = id
	}
}

func WithSessionAccountID(accountID uuid.UUID) func(*auth.Session) {
	return func(s *auth.Session) {
		s.AccountID = accountID
	}
}

func WithSessionToken(token string) func(*auth.Session) {
	return func(s *auth.Session) {
		s.Token = token
	}
}

func WithUnhashedToken(token string) func(*auth.Session) {
	return func(s *auth.Session) {
		s.UnhashedToken = token
	}
}

func WithRevokedAt(t time.Time) func(*auth.Session) {
	return func(s *auth.Session) {
		s.RevokedAt = &t
	}
}

func WithCreatedAt(t time.Time) func(*auth.Session) {
	return func(s *auth.Session) {
		s.CreatedAt = t
	}
}

func WithRotatedAt(t time.Time) func(*auth.Session) {
	return func(s *auth.Session) {
		s.RotatedAt = t
	}
}

func NewTestAuthUser(opts ...func(*auth.AuthUser)) *auth.AuthUser {
	user := &auth.AuthUser{
		AccountID: uuid.New(),
		Email:     "test@example.com",
		IsAdmin:   false,
		SessionID: uuid.New(),
	}

	for _, opt := range opts {
		opt(user)
	}

	return user
}

func WithAuthUserAccountID(id uuid.UUID) func(*auth.AuthUser) {
	return func(u *auth.AuthUser) {
		u.AccountID = id
	}
}

func WithAuthUserIsAdmin(isAdmin bool) func(*auth.AuthUser) {
	return func(u *auth.AuthUser) {
		u.IsAdmin = isAdmin
	}
}

func NewTestAuthContext(opts ...func(*auth.AuthContext)) *auth.AuthContext {
	now := time.Now()
	ctx := &auth.AuthContext{
		SessionID:      uuid.New(),
		SessionToken:   "hashed_token",
		AccountID:      uuid.New(),
		AccountEmail:   "test@example.com",
		AccountIsAdmin: false,
		AccountStatus:  "active",
		SessionCreated: now,
		SessionRotated: now,
	}

	for _, opt := range opts {
		opt(ctx)
	}

	return ctx
}

func WithAuthContextAccountID(id uuid.UUID) func(*auth.AuthContext) {
	return func(c *auth.AuthContext) {
		c.AccountID = id
	}
}

func WithAuthContextSessionRevoked(t time.Time) func(*auth.AuthContext) {
	return func(c *auth.AuthContext) {
		c.SessionRevoked = &t
	}
}

func WithAuthContextSessionCreated(t time.Time) func(*auth.AuthContext) {
	return func(c *auth.AuthContext) {
		c.SessionCreated = t
	}
}

func WithAuthContextSessionRotated(t time.Time) func(*auth.AuthContext) {
	return func(c *auth.AuthContext) {
		c.SessionRotated = t
	}
}

func StringPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}
