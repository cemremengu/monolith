package testutil

import (
	"context"
	"time"

	"monolith/internal/config"
	"monolith/internal/database"
	"monolith/internal/service/account"
	"monolith/internal/service/auth"
	"monolith/internal/service/login"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
)

func NewMockDB() (pgxmock.PgxPoolIface, *database.DB) {
	mock, _ := pgxmock.NewPool()
	db := &database.DB{Pool: mock}
	return mock, db
}

func NewTestSecurityConfig() config.SecurityConfig {
	return config.SecurityConfig{
		SecretKey:                            "test-secret-key",
		TokenSecretKey:                       "test-token-secret-key",
		LoginMaximumLifetimeDuration:         30 * 24 * time.Hour,
		LoginMaximumInactiveLifetimeDuration: 7 * 24 * time.Hour,
		LoginCookieName:                      "session_token",
		TokenRotationIntervalMinutes:         10,
	}
}

type MockAccountService struct {
	ValidatePasswordFn  func(hashedPassword, password string) error
	UserExistsFn        func(ctx context.Context, email, username string) (bool, error)
	RegisterFn          func(ctx context.Context, req account.RegisterRequest) (*account.Account, error)
	GetAccountByIDFn    func(ctx context.Context, accountID uuid.UUID) (*account.Account, error)
	GetAccountByLoginFn func(ctx context.Context, login string) (*account.Account, error)
	UpdatePreferencesFn func(ctx context.Context, accountID uuid.UUID, req account.UpdatePreferencesRequest) (*account.Account, error)
	ChangePasswordFn    func(ctx context.Context, accountID uuid.UUID, req account.ChangePasswordRequest) error
	UpdateLastSeenFn    func(ctx context.Context, accountID uuid.UUID) error
	GetAccountsFn       func(ctx context.Context) ([]account.Account, error)
	GetAccountFn        func(ctx context.Context, id uuid.UUID) (*account.Account, error)
	CreateAccountFn     func(ctx context.Context, req account.CreateAccountRequest) (*account.Account, error)
	UpdateAccountFn     func(ctx context.Context, id uuid.UUID, req account.UpdateAccountRequest) (*account.Account, error)
	DisableAccountFn    func(ctx context.Context, id uuid.UUID) error
	EnableAccountFn     func(ctx context.Context, id uuid.UUID) error
	DeleteAccountFn     func(ctx context.Context, id uuid.UUID) error
	InviteUsersFn       func(ctx context.Context, req account.InviteUsersRequest) (*account.InviteUsersResponse, error)
}

func (m *MockAccountService) ValidatePassword(hashedPassword, password string) error {
	if m.ValidatePasswordFn != nil {
		return m.ValidatePasswordFn(hashedPassword, password)
	}
	return nil
}

func (m *MockAccountService) UserExists(ctx context.Context, email, username string) (bool, error) {
	if m.UserExistsFn != nil {
		return m.UserExistsFn(ctx, email, username)
	}
	return false, nil
}

func (m *MockAccountService) Register(ctx context.Context, req account.RegisterRequest) (*account.Account, error) {
	if m.RegisterFn != nil {
		return m.RegisterFn(ctx, req)
	}
	return NewTestAccount(), nil
}

func (m *MockAccountService) GetAccountByID(ctx context.Context, accountID uuid.UUID) (*account.Account, error) {
	if m.GetAccountByIDFn != nil {
		return m.GetAccountByIDFn(ctx, accountID)
	}
	return NewTestAccount(WithAccountID(accountID)), nil
}

func (m *MockAccountService) GetAccountByLogin(ctx context.Context, loginStr string) (*account.Account, error) {
	if m.GetAccountByLoginFn != nil {
		return m.GetAccountByLoginFn(ctx, loginStr)
	}
	return NewTestAccount(), nil
}

func (m *MockAccountService) UpdatePreferences(ctx context.Context, accountID uuid.UUID, req account.UpdatePreferencesRequest) (*account.Account, error) {
	if m.UpdatePreferencesFn != nil {
		return m.UpdatePreferencesFn(ctx, accountID, req)
	}
	return NewTestAccount(WithAccountID(accountID)), nil
}

func (m *MockAccountService) ChangePassword(ctx context.Context, accountID uuid.UUID, req account.ChangePasswordRequest) error {
	if m.ChangePasswordFn != nil {
		return m.ChangePasswordFn(ctx, accountID, req)
	}
	return nil
}

func (m *MockAccountService) UpdateLastSeen(ctx context.Context, accountID uuid.UUID) error {
	if m.UpdateLastSeenFn != nil {
		return m.UpdateLastSeenFn(ctx, accountID)
	}
	return nil
}

func (m *MockAccountService) GetAccounts(ctx context.Context) ([]account.Account, error) {
	if m.GetAccountsFn != nil {
		return m.GetAccountsFn(ctx)
	}
	return []account.Account{*NewTestAccount()}, nil
}

func (m *MockAccountService) GetAccount(ctx context.Context, id uuid.UUID) (*account.Account, error) {
	if m.GetAccountFn != nil {
		return m.GetAccountFn(ctx, id)
	}
	return NewTestAccount(WithAccountID(id)), nil
}

func (m *MockAccountService) CreateAccount(ctx context.Context, req account.CreateAccountRequest) (*account.Account, error) {
	if m.CreateAccountFn != nil {
		return m.CreateAccountFn(ctx, req)
	}
	return NewTestAccount(), nil
}

func (m *MockAccountService) UpdateAccount(ctx context.Context, id uuid.UUID, req account.UpdateAccountRequest) (*account.Account, error) {
	if m.UpdateAccountFn != nil {
		return m.UpdateAccountFn(ctx, id, req)
	}
	return NewTestAccount(WithAccountID(id)), nil
}

func (m *MockAccountService) DisableAccount(ctx context.Context, id uuid.UUID) error {
	if m.DisableAccountFn != nil {
		return m.DisableAccountFn(ctx, id)
	}
	return nil
}

func (m *MockAccountService) EnableAccount(ctx context.Context, id uuid.UUID) error {
	if m.EnableAccountFn != nil {
		return m.EnableAccountFn(ctx, id)
	}
	return nil
}

func (m *MockAccountService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	if m.DeleteAccountFn != nil {
		return m.DeleteAccountFn(ctx, id)
	}
	return nil
}

func (m *MockAccountService) InviteUsers(ctx context.Context, req account.InviteUsersRequest) (*account.InviteUsersResponse, error) {
	if m.InviteUsersFn != nil {
		return m.InviteUsersFn(ctx, req)
	}
	return &account.InviteUsersResponse{Success: []account.Account{*NewTestAccount()}}, nil
}

type MockAuthService struct {
	CreateSessionFn         func(ctx context.Context, req *auth.CreateSessionRequest) (*auth.Session, error)
	GetAuthContextByTokenFn func(ctx context.Context, unhashedToken string) (*auth.AuthContext, error)
	GetSessionByTokenFn     func(ctx context.Context, unhashedToken string) (*auth.Session, error)
	RotateSessionFn         func(ctx context.Context, req *auth.RotateSessionRequest) (*auth.Session, error)
	RevokeSessionFn         func(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
	GetUserSessionsFn       func(ctx context.Context, userID uuid.UUID) ([]auth.UserSession, error)
	RevokeAllUserSessionsFn func(ctx context.Context, accountID uuid.UUID) error
}

func (m *MockAuthService) CreateSession(ctx context.Context, req *auth.CreateSessionRequest) (*auth.Session, error) {
	if m.CreateSessionFn != nil {
		return m.CreateSessionFn(ctx, req)
	}
	return NewTestSession(WithSessionAccountID(req.AccountID)), nil
}

func (m *MockAuthService) GetAuthContextByToken(ctx context.Context, unhashedToken string) (*auth.AuthContext, error) {
	if m.GetAuthContextByTokenFn != nil {
		return m.GetAuthContextByTokenFn(ctx, unhashedToken)
	}
	return NewTestAuthContext(), nil
}

func (m *MockAuthService) GetSessionByToken(ctx context.Context, unhashedToken string) (*auth.Session, error) {
	if m.GetSessionByTokenFn != nil {
		return m.GetSessionByTokenFn(ctx, unhashedToken)
	}
	return NewTestSession(), nil
}

func (m *MockAuthService) RotateSession(ctx context.Context, req *auth.RotateSessionRequest) (*auth.Session, error) {
	if m.RotateSessionFn != nil {
		return m.RotateSessionFn(ctx, req)
	}
	return NewTestSession(), nil
}

func (m *MockAuthService) RevokeSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	if m.RevokeSessionFn != nil {
		return m.RevokeSessionFn(ctx, userID, sessionID)
	}
	return nil
}

func (m *MockAuthService) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]auth.UserSession, error) {
	if m.GetUserSessionsFn != nil {
		return m.GetUserSessionsFn(ctx, userID)
	}
	return []auth.UserSession{}, nil
}

func (m *MockAuthService) RevokeAllUserSessions(ctx context.Context, accountID uuid.UUID) error {
	if m.RevokeAllUserSessionsFn != nil {
		return m.RevokeAllUserSessionsFn(ctx, accountID)
	}
	return nil
}

type MockLoginService struct {
	LoginFn func(ctx context.Context, req login.UserLoginRequest) (*account.Account, error)
}

func (m *MockLoginService) Login(ctx context.Context, req login.UserLoginRequest) (*account.Account, error) {
	if m.LoginFn != nil {
		return m.LoginFn(ctx, req)
	}
	return NewTestAccount(), nil
}
