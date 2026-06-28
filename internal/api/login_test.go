package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"monolith/internal/account"
	"monolith/internal/auth"
	"monolith/internal/config"
	"monolith/internal/database"
	"monolith/internal/login"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func newTestSecurityConfig() config.SecurityConfig {
	return config.SecurityConfig{
		SecretKey:                            "test-secret-key",
		LoginMaximumLifetimeDuration:         30 * 24 * time.Hour,
		LoginMaximumInactiveLifetimeDuration: 7 * 24 * time.Hour,
		LoginCookieName:                      "session_token",
		TokenRotationIntervalMinutes:         10,
	}
}

func TestAuthHandler_Login(t *testing.T) {
	accountID := uuid.New()
	sessionID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	now := time.Now()
	cfg := newTestSecurityConfig()

	tests := []struct {
		name       string
		body       map[string]any
		setupMock  func(mock pgxmock.PgxPoolIface)
		wantStatus int
		wantCookie bool
	}{
		{
			name: "valid credentials",
			body: map[string]any{
				"login":    "test@example.com",
				"password": "password123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				accountRows := pgxmock.NewRows([]string{
					"id", "username", "email", "name", "password", "is_admin", "language", "theme", "timezone",
					"last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(
					accountID, "testuser", "test@example.com", new("Test User"),
					string(hashedPassword), false, new("en"), new("light"),
					new("UTC"), nil, "active", now, now,
				)
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("test@example.com").
					WillReturnRows(accountRows)

				mock.ExpectExec(`UPDATE account SET last_seen_at = NOW\(\) WHERE id = \$1`).
					WithArgs(accountID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				sessionRows := pgxmock.NewRows([]string{
					"id", "token", "prev_token", "account_id", "user_agent", "client_ip",
					"token_seen", "seen_at", "created_at", "rotated_at", "revoked_at",
				}).AddRow(
					sessionID, "hashed_token", new("hashed_token"), accountID,
					"", "", false, (*time.Time)(nil), now, now, (*time.Time)(nil),
				)
				mock.ExpectQuery(`INSERT INTO auth_session`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), accountID, pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(sessionRows)
			},
			wantStatus: http.StatusOK,
			wantCookie: true,
		},
		{
			name: "invalid credentials - wrong password",
			body: map[string]any{
				"login":    "test@example.com",
				"password": "wrongpassword",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				accountRows := pgxmock.NewRows([]string{
					"id", "username", "email", "name", "password", "is_admin", "language", "theme", "timezone",
					"last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(
					accountID, "testuser", "test@example.com", new("Test User"),
					string(hashedPassword), false, new("en"), new("light"),
					new("UTC"), nil, "active", now, now,
				)
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("test@example.com").
					WillReturnRows(accountRows)
			},
			wantStatus: http.StatusUnauthorized,
			wantCookie: false,
		},
		{
			name: "invalid credentials - user not found",
			body: map[string]any{
				"login":    "notfound@example.com",
				"password": "password123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("notfound@example.com").
					WillReturnError(database.ErrNoRows)
			},
			wantStatus: http.StatusUnauthorized,
			wantCookie: false,
		},
		{
			name:       "malformed JSON",
			body:       nil,
			setupMock:  func(mock pgxmock.PgxPoolIface) {},
			wantStatus: http.StatusBadRequest,
			wantCookie: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			accountService := account.NewService(db)
			loginService := login.NewService(db, accountService)
			authService := auth.NewService(db, cfg)
			handler := NewAuthHandler(loginService, authService)

			e := echo.New()
			var req *http.Request
			if tt.body != nil {
				jsonBody, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(jsonBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			} else {
				req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader([]byte("invalid json")))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Login(c)

			if tt.wantStatus >= 400 {
				require.Error(t, err)
				httpErr := &echo.HTTPError{}
				ok := errors.As(err, &httpErr)
				require.True(t, ok)
				assert.Equal(t, tt.wantStatus, httpErr.Code)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
			}

			if tt.wantCookie {
				cookies := rec.Result().Cookies()
				var hasSessionCookie bool
				for _, cookie := range cookies {
					if cookie.Name == cfg.LoginCookieName {
						hasSessionCookie = true
						break
					}
				}
				assert.True(t, hasSessionCookie, "should set session cookie")
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	cfg := newTestSecurityConfig()
	hashedToken := auth.HashTokenForTest("valid_token", cfg.SecretKey)

	tests := []struct {
		name        string
		cookieValue string
		setupMock   func(mock pgxmock.PgxPoolIface)
		wantStatus  int
		wantErr     bool
	}{
		{
			name:        "revokes session and clears cookies",
			cookieValue: "valid_token",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE auth_session SET revoked_at = NOW\(\) WHERE \(token = \$1 OR prev_token = \$2\) AND revoked_at IS NULL`).
					WithArgs(hashedToken, hashedToken).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "clears cookies when session cookie is missing",
			setupMock:  func(mock pgxmock.PgxPoolIface) {},
			wantStatus: http.StatusOK,
		},
		{
			name:        "returns error when revocation fails",
			cookieValue: "valid_token",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE auth_session SET revoked_at = NOW\(\) WHERE \(token = \$1 OR prev_token = \$2\) AND revoked_at IS NULL`).
					WithArgs(hashedToken, hashedToken).
					WillReturnError(assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			accountService := account.NewService(db)
			loginService := login.NewService(db, accountService)
			authService := auth.NewService(db, cfg)
			handler := NewAuthHandler(loginService, authService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{Name: cfg.LoginCookieName, Value: tt.cookieValue})
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Logout(c)
			if tt.wantErr {
				require.Error(t, err)
				httpErr := &echo.HTTPError{}
				ok := errors.As(err, &httpErr)
				require.True(t, ok)
				assert.Equal(t, tt.wantStatus, httpErr.Code)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
			}

			cookies := rec.Result().Cookies()
			var hasClearedSessionCookie bool
			for _, cookie := range cookies {
				if cookie.Name == cfg.LoginCookieName {
					hasClearedSessionCookie = true
					assert.Equal(t, -1, cookie.MaxAge, "session cookie should be cleared")
				}
			}
			assert.True(t, hasClearedSessionCookie, "session cookie should be cleared")
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
