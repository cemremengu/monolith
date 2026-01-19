package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"monolith/internal/config"
	"monolith/internal/database"
	"monolith/internal/service/auth"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSecurityConfig() config.SecurityConfig {
	return config.SecurityConfig{
		SecretKey:                            "test-secret-key",
		TokenSecretKey:                       "test-token-secret-key",
		LoginMaximumLifetimeDuration:         30 * 24 * time.Hour,
		LoginMaximumInactiveLifetimeDuration: 7 * 24 * time.Hour,
		LoginCookieName:                      "session_token",
		TokenRotationIntervalMinutes:         10,
	}
}

func TestSessionAuth(t *testing.T) {
	cfg := newTestSecurityConfig()
	sessionID := uuid.New()
	accountID := uuid.New()
	now := time.Now()

	hashedToken := auth.HashTokenForTest("valid_token", cfg.TokenSecretKey)

	tests := []struct {
		name         string
		cookie       *http.Cookie
		setupMock    func(mock pgxmock.PgxPoolIface)
		wantStatus   int
		wantUserSet  bool
		handlerCalls int
	}{
		{
			name: "valid cookie",
			cookie: &http.Cookie{
				Name:  cfg.LoginCookieName,
				Value: "valid_token",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"session_id", "session_token", "account_id", "account_email",
					"account_is_admin", "account_status", "session_created", "session_rotated", "session_revoked",
				}).AddRow(
					sessionID, hashedToken, accountID, "test@example.com",
					false, "active", now, now, (*time.Time)(nil),
				)
				mock.ExpectQuery(`SELECT .+ FROM auth_session s INNER JOIN account a ON s.account_id = a.id`).
					WithArgs(hashedToken, hashedToken, auth.AccountStatusActive).
					WillReturnRows(rows)
			},
			wantStatus:   http.StatusOK,
			wantUserSet:  true,
			handlerCalls: 1,
		},
		{
			name:         "missing cookie",
			cookie:       nil,
			setupMock:    func(mock pgxmock.PgxPoolIface) {},
			wantStatus:   http.StatusUnauthorized,
			wantUserSet:  false,
			handlerCalls: 0,
		},
		{
			name: "invalid token",
			cookie: &http.Cookie{
				Name:  cfg.LoginCookieName,
				Value: "invalid_token",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT .+ FROM auth_session s INNER JOIN account a ON s.account_id = a.id`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), auth.AccountStatusActive).
					WillReturnError(database.ErrNoRows)
			},
			wantStatus:   http.StatusUnauthorized,
			wantUserSet:  false,
			handlerCalls: 0,
		},
		{
			name: "expired session",
			cookie: &http.Cookie{
				Name:  cfg.LoginCookieName,
				Value: "valid_token",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				oldCreated := time.Now().Add(-31 * 24 * time.Hour)
				rows := pgxmock.NewRows([]string{
					"session_id", "session_token", "account_id", "account_email",
					"account_is_admin", "account_status", "session_created", "session_rotated", "session_revoked",
				}).AddRow(
					sessionID, hashedToken, accountID, "test@example.com",
					false, "active", oldCreated, now, (*time.Time)(nil),
				)
				mock.ExpectQuery(`SELECT .+ FROM auth_session s INNER JOIN account a ON s.account_id = a.id`).
					WithArgs(hashedToken, hashedToken, auth.AccountStatusActive).
					WillReturnRows(rows)
			},
			wantStatus:   http.StatusUnauthorized,
			wantUserSet:  false,
			handlerCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			authService := auth.NewService(db, cfg)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handlerCalled := 0
			handler := func(c echo.Context) error {
				handlerCalled++
				return c.String(http.StatusOK, "OK")
			}

			middleware := SessionAuth(authService, cfg)
			err := middleware(handler)(c)

			if tt.wantStatus == http.StatusOK {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.handlerCalls, handlerCalled)

			if tt.wantUserSet {
				user := c.Get("user")
				assert.NotNil(t, user)
				authUser, ok := user.(*auth.AuthUser)
				require.True(t, ok)
				assert.Equal(t, accountID, authUser.AccountID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
