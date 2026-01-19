package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"monolith/internal/database"
	"monolith/internal/service/account"
	"monolith/internal/service/auth"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockValidator struct{}

func (v *mockValidator) Validate(i any) error {
	return nil
}

func TestAccountHandler_Profile(t *testing.T) {
	accountID := uuid.New()
	now := time.Now()

	tests := []struct {
		name       string
		user       *auth.AuthUser
		setupMock  func(mock pgxmock.PgxPoolIface)
		wantStatus int
	}{
		{
			name: "authenticated user",
			user: &auth.AuthUser{
				AccountID: accountID,
				Email:     "test@example.com",
				IsAdmin:   false,
				SessionID: uuid.New(),
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "username", "email", "name", "is_admin", "language", "theme", "timezone",
					"last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(
					accountID, "testuser", "test@example.com", stringPtr("Test User"),
					false, stringPtr("en"), stringPtr("light"),
					stringPtr("UTC"), nil, "active", now, now,
				)
				mock.ExpectQuery(`SELECT .+ FROM account WHERE id = \$1 AND status = 'active'`).
					WithArgs(accountID).
					WillReturnRows(rows)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "unauthenticated",
			user:       nil,
			setupMock:  func(mock pgxmock.PgxPoolIface) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "account not found",
			user: &auth.AuthUser{
				AccountID: uuid.New(),
				Email:     "test@example.com",
				IsAdmin:   false,
				SessionID: uuid.New(),
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT .+ FROM account WHERE id = \$1 AND status = 'active'`).
					WithArgs(pgxmock.AnyArg()).
					WillReturnError(database.ErrNoRows)
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			accountService := account.NewService(db)
			handler := NewAccountHandler(accountService)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/account/profile", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.user != nil {
				c.Set("user", tt.user)
			}

			err := handler.Profile(c)

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

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAccountHandler_Register(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]any
		setupMock  func(mock pgxmock.PgxPoolIface)
		wantStatus int
	}{
		{
			name: "password too short",
			body: map[string]any{
				"username": "newuser",
				"email":    "newuser@example.com",
				"password": "short",
				"name":     "New User",
			},
			setupMock:  func(mock pgxmock.PgxPoolIface) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "user already exists",
			body: map[string]any{
				"username": "existinguser",
				"email":    "existing@example.com",
				"password": "password123",
				"name":     "Existing User",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(uuid.New())
				mock.ExpectQuery(`SELECT id FROM account WHERE email = \$1 OR username = \$2`).
					WithArgs("existing@example.com", "existinguser").
					WillReturnRows(rows)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:       "malformed JSON",
			body:       nil,
			setupMock:  func(mock pgxmock.PgxPoolIface) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			accountService := account.NewService(db)
			handler := NewAccountHandler(accountService)

			e := echo.New()
			e.Validator = &mockValidator{}

			var req *http.Request
			if tt.body != nil {
				jsonBody, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(http.MethodPost, "/api/account/register", bytes.NewReader(jsonBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			} else {
				req = httptest.NewRequest(http.MethodPost, "/api/account/register", bytes.NewReader([]byte("invalid")))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Register(c)

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

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAccountHandler_UpdatePreferences(t *testing.T) {
	accountID := uuid.New()
	now := time.Now()

	tests := []struct {
		name       string
		user       *auth.AuthUser
		body       map[string]any
		setupMock  func(mock pgxmock.PgxPoolIface)
		wantStatus int
	}{
		{
			name: "successful update",
			user: &auth.AuthUser{
				AccountID: accountID,
				Email:     "test@example.com",
				IsAdmin:   false,
				SessionID: uuid.New(),
			},
			body: map[string]any{
				"language": "fr",
				"theme":    "dark",
				"timezone": "Europe/Paris",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "username", "email", "name", "is_admin", "language", "theme", "timezone",
					"last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(
					accountID, "testuser", "test@example.com", stringPtr("Test User"),
					false, stringPtr("fr"), stringPtr("dark"),
					stringPtr("Europe/Paris"), nil, "active", now, now,
				)
				mock.ExpectQuery(`UPDATE account SET language = \$1, theme = \$2, timezone = \$3, updated_at = NOW\(\) WHERE id = \$4 AND status = 'active' RETURNING`).
					WithArgs(stringPtr("fr"), stringPtr("dark"), stringPtr("Europe/Paris"), accountID).
					WillReturnRows(rows)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "unauthenticated",
			user:       nil,
			body:       map[string]any{"language": "fr"},
			setupMock:  func(mock pgxmock.PgxPoolIface) {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			accountService := account.NewService(db)
			handler := NewAccountHandler(accountService)

			e := echo.New()
			e.Validator = &mockValidator{}

			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/api/account/preferences", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.user != nil {
				c.Set("user", tt.user)
			}

			err := handler.UpdatePreferences(c)

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

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAccountHandler_ChangePassword(t *testing.T) {
	accountID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("currentpass123"), 12)

	tests := []struct {
		name       string
		user       *auth.AuthUser
		body       map[string]any
		setupMock  func(mock pgxmock.PgxPoolIface)
		wantStatus int
	}{
		{
			name: "successful password change",
			user: &auth.AuthUser{
				AccountID: accountID,
				Email:     "test@example.com",
				IsAdmin:   false,
				SessionID: uuid.New(),
			},
			body: map[string]any{
				"currentPassword": "currentpass123",
				"newPassword":     "newpassword123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "password"}).
					AddRow(accountID, string(hashedPassword))
				mock.ExpectQuery(`SELECT id, password FROM account WHERE id = \$1 AND status = 'active'`).
					WithArgs(accountID).
					WillReturnRows(rows)

				mock.ExpectExec(`UPDATE account SET password = \$1, updated_at = NOW\(\) WHERE id = \$2`).
					WithArgs(pgxmock.AnyArg(), accountID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "unauthenticated",
			user:       nil,
			body:       map[string]any{"currentPassword": "old", "newPassword": "newpassword123"},
			setupMock:  func(mock pgxmock.PgxPoolIface) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong current password",
			user: &auth.AuthUser{
				AccountID: accountID,
				Email:     "test@example.com",
				IsAdmin:   false,
				SessionID: uuid.New(),
			},
			body: map[string]any{
				"currentPassword": "wrongpassword",
				"newPassword":     "newpassword123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "password"}).
					AddRow(accountID, string(hashedPassword))
				mock.ExpectQuery(`SELECT id, password FROM account WHERE id = \$1 AND status = 'active'`).
					WithArgs(accountID).
					WillReturnRows(rows)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "new password too short",
			user: &auth.AuthUser{
				AccountID: accountID,
				Email:     "test@example.com",
				IsAdmin:   false,
				SessionID: uuid.New(),
			},
			body: map[string]any{
				"currentPassword": "currentpass123",
				"newPassword":     "short",
			},
			setupMock:  func(mock pgxmock.PgxPoolIface) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			accountService := account.NewService(db)
			handler := NewAccountHandler(accountService)

			e := echo.New()
			e.Validator = &mockValidator{}

			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/account/change-password", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.user != nil {
				c.Set("user", tt.user)
			}

			err := handler.ChangePassword(c)

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

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
