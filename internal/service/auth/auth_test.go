package auth

import (
	"context"
	"testing"
	"time"

	"monolith/internal/config"
	"monolith/internal/database"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func newTestService(mock pgxmock.PgxPoolIface) *Service {
	db := &database.DB{Pool: mock}
	cfg := newTestSecurityConfig()
	return NewService(db, cfg)
}

func stringPtr(s string) *string {
	return &s
}

func TestService_CreateSession(t *testing.T) {
	accountID := uuid.New()
	sessionID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		req       *CreateSessionRequest
		setupMock func(mock pgxmock.PgxPoolIface)
		wantErr   bool
	}{
		{
			name: "successful session creation",
			req: &CreateSessionRequest{
				AccountID: accountID,
				ClientIP:  "127.0.0.1",
				UserAgent: "Mozilla/5.0 Test Browser",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "token", "prev_token", "account_id", "user_agent", "client_ip",
					"token_seen", "seen_at", "created_at", "rotated_at", "revoked_at",
				}).AddRow(
					sessionID, "hashed_token", stringPtr("hashed_token"), accountID,
					"Mozilla/5.0 Test Browser", "127.0.0.1",
					false, (*time.Time)(nil), now, now, (*time.Time)(nil),
				)
				mock.ExpectQuery(`INSERT INTO auth_session`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), accountID, "Mozilla/5.0 Test Browser", "127.0.0.1").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "database error",
			req: &CreateSessionRequest{
				AccountID: accountID,
				ClientIP:  "127.0.0.1",
				UserAgent: "Mozilla/5.0 Test Browser",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`INSERT INTO auth_session`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), accountID, "Mozilla/5.0 Test Browser", "127.0.0.1").
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			s := newTestService(mock)

			session, err := s.CreateSession(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, session)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, session)
				assert.NotEmpty(t, session.UnhashedToken)
				assert.Equal(t, accountID, session.AccountID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_GetAuthContextByToken(t *testing.T) {
	sessionID := uuid.New()
	accountID := uuid.New()
	now := time.Now()

	cfg := newTestSecurityConfig()
	hashedToken := hashToken("valid_token", cfg.SecretKey)

	tests := []struct {
		name      string
		token     string
		setupMock func(mock pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:  "valid token",
			token: "valid_token",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"session_id", "session_token", "account_id", "account_email",
					"account_is_admin", "account_status", "session_created", "session_rotated", "session_revoked",
				}).AddRow(
					sessionID, hashedToken, accountID, "test@example.com",
					false, "active", now, now, (*time.Time)(nil),
				)
				mock.ExpectQuery(`SELECT .+ FROM auth_session s INNER JOIN account a ON s.account_id = a.id`).
					WithArgs(hashedToken, hashedToken, AccountStatusActive).
					WillReturnRows(rows)
			},
			wantErr: nil,
		},
		{
			name:  "session not found",
			token: "invalid_token",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT .+ FROM auth_session s INNER JOIN account a ON s.account_id = a.id`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), AccountStatusActive).
					WillReturnError(database.ErrNoRows)
			},
			wantErr: ErrSessionNotFound,
		},
		{
			name:  "session revoked",
			token: "valid_token",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				revokedAt := time.Now()
				rows := pgxmock.NewRows([]string{
					"session_id", "session_token", "account_id", "account_email",
					"account_is_admin", "account_status", "session_created", "session_rotated", "session_revoked",
				}).AddRow(
					sessionID, hashedToken, accountID, "test@example.com",
					false, "active", now, now, &revokedAt,
				)
				mock.ExpectQuery(`SELECT .+ FROM auth_session s INNER JOIN account a ON s.account_id = a.id`).
					WithArgs(hashedToken, hashedToken, AccountStatusActive).
					WillReturnRows(rows)
			},
			wantErr: ErrSessionRevoked,
		},
		{
			name:  "session expired by created_at",
			token: "valid_token",
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
					WithArgs(hashedToken, hashedToken, AccountStatusActive).
					WillReturnRows(rows)
			},
			wantErr: ErrSessionExpired,
		},
		{
			name:  "session expired by rotated_at (inactive)",
			token: "valid_token",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				oldRotated := time.Now().Add(-8 * 24 * time.Hour)
				rows := pgxmock.NewRows([]string{
					"session_id", "session_token", "account_id", "account_email",
					"account_is_admin", "account_status", "session_created", "session_rotated", "session_revoked",
				}).AddRow(
					sessionID, hashedToken, accountID, "test@example.com",
					false, "active", now, oldRotated, (*time.Time)(nil),
				)
				mock.ExpectQuery(`SELECT .+ FROM auth_session s INNER JOIN account a ON s.account_id = a.id`).
					WithArgs(hashedToken, hashedToken, AccountStatusActive).
					WillReturnRows(rows)
			},
			wantErr: ErrSessionExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			s := newTestService(mock)

			authCtx, err := s.GetAuthContextByToken(context.Background(), tt.token)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, authCtx)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, authCtx)
				assert.Equal(t, accountID, authCtx.AccountID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_GetSessionByToken(t *testing.T) {
	sessionID := uuid.New()
	accountID := uuid.New()
	now := time.Now()

	cfg := newTestSecurityConfig()
	hashedToken := hashToken("valid_token", cfg.SecretKey)

	tests := []struct {
		name      string
		token     string
		setupMock func(mock pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:  "valid token",
			token: "valid_token",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "token", "account_id", "user_agent", "client_ip",
					"created_at", "rotated_at", "revoked_at",
				}).AddRow(
					sessionID, hashedToken, accountID,
					"Mozilla/5.0 Test Browser", "127.0.0.1",
					now, now, (*time.Time)(nil),
				)
				mock.ExpectQuery(`SELECT .+ FROM auth_session WHERE token = \$1 OR prev_token = \$2`).
					WithArgs(hashedToken, hashedToken).
					WillReturnRows(rows)
			},
			wantErr: nil,
		},
		{
			name:  "session not found",
			token: "invalid_token",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT .+ FROM auth_session WHERE token = \$1 OR prev_token = \$2`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(database.ErrNoRows)
			},
			wantErr: ErrSessionNotFound,
		},
		{
			name:  "session revoked",
			token: "valid_token",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				revokedAt := time.Now()
				rows := pgxmock.NewRows([]string{
					"id", "token", "account_id", "user_agent", "client_ip",
					"created_at", "rotated_at", "revoked_at",
				}).AddRow(
					sessionID, hashedToken, accountID,
					"Mozilla/5.0 Test Browser", "127.0.0.1",
					now, now, &revokedAt,
				)
				mock.ExpectQuery(`SELECT .+ FROM auth_session WHERE token = \$1 OR prev_token = \$2`).
					WithArgs(hashedToken, hashedToken).
					WillReturnRows(rows)
			},
			wantErr: ErrSessionRevoked,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			s := newTestService(mock)

			session, err := s.GetSessionByToken(context.Background(), tt.token)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_RotateSession(t *testing.T) {
	sessionID := uuid.New()
	accountID := uuid.New()
	now := time.Now()

	cfg := newTestSecurityConfig()
	hashedToken := hashToken("old_token", cfg.SecretKey)

	tests := []struct {
		name      string
		req       *RotateSessionRequest
		setupMock func(mock pgxmock.PgxPoolIface)
		wantErr   bool
	}{
		{
			name: "successful rotation",
			req: &RotateSessionRequest{
				UnhashedToken: "old_token",
				ClientIP:      "127.0.0.1",
				UserAgent:     "Mozilla/5.0 Test Browser",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				getRows := pgxmock.NewRows([]string{
					"id", "token", "account_id", "user_agent", "client_ip",
					"created_at", "rotated_at", "revoked_at",
				}).AddRow(
					sessionID, hashedToken, accountID,
					"Mozilla/5.0 Test Browser", "127.0.0.1",
					now, now, (*time.Time)(nil),
				)
				mock.ExpectQuery(`SELECT .+ FROM auth_session WHERE token = \$1 OR prev_token = \$2`).
					WithArgs(hashedToken, hashedToken).
					WillReturnRows(getRows)

				updateRows := pgxmock.NewRows([]string{
					"id", "token", "prev_token", "account_id", "user_agent", "client_ip",
					"token_seen", "seen_at", "created_at", "rotated_at", "revoked_at",
				}).AddRow(
					sessionID, "new_hashed_token", stringPtr(hashedToken), accountID,
					"Mozilla/5.0 Test Browser", "127.0.0.1",
					false, (*time.Time)(nil), now, time.Now(), (*time.Time)(nil),
				)
				mock.ExpectQuery(`UPDATE auth_session SET token = \$1, prev_token = \$2, rotated_at = NOW\(\), token_seen = FALSE, seen_at = NULL WHERE id = \$3 RETURNING`).
					WithArgs(pgxmock.AnyArg(), hashedToken, sessionID).
					WillReturnRows(updateRows)
			},
			wantErr: false,
		},
		{
			name: "session not found",
			req: &RotateSessionRequest{
				UnhashedToken: "invalid_token",
				ClientIP:      "127.0.0.1",
				UserAgent:     "Mozilla/5.0 Test Browser",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT .+ FROM auth_session WHERE token = \$1 OR prev_token = \$2`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(database.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			s := newTestService(mock)

			session, err := s.RotateSession(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.NotEmpty(t, session.UnhashedToken)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_RevokeSession(t *testing.T) {
	userID := uuid.New()
	sessionID := uuid.New()

	tests := []struct {
		name      string
		userID    uuid.UUID
		sessionID uuid.UUID
		setupMock func(mock pgxmock.PgxPoolIface)
		wantErr   bool
	}{
		{
			name:      "successful revocation",
			userID:    userID,
			sessionID: sessionID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE auth_session SET revoked_at = NOW\(\) WHERE id = \$1 and account_id = \$2`).
					WithArgs(sessionID, userID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: false,
		},
		{
			name:      "database error",
			userID:    userID,
			sessionID: sessionID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE auth_session SET revoked_at = NOW\(\) WHERE id = \$1 and account_id = \$2`).
					WithArgs(sessionID, userID).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			s := newTestService(mock)

			err := s.RevokeSession(context.Background(), tt.userID, tt.sessionID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_GetUserSessions(t *testing.T) {
	userID := uuid.New()
	sessionID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		userID    uuid.UUID
		setupMock func(mock pgxmock.PgxPoolIface)
		wantCount int
		wantErr   bool
	}{
		{
			name:   "returns sessions",
			userID: userID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "token", "account_id", "user_agent", "client_ip",
					"created_at", "rotated_at", "revoked_at",
				}).AddRow(
					sessionID, "token1", userID,
					"Mozilla/5.0", "127.0.0.1",
					now, now, (*time.Time)(nil),
				).AddRow(
					uuid.New(), "token2", userID,
					"Safari/1.0", "192.168.1.1",
					now, now, (*time.Time)(nil),
				)
				mock.ExpectQuery(`SELECT .+ FROM auth_session WHERE account_id = \$1 AND revoked_at IS NULL ORDER BY rotated_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:   "no sessions",
			userID: userID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "token", "account_id", "user_agent", "client_ip",
					"created_at", "rotated_at", "revoked_at",
				})
				mock.ExpectQuery(`SELECT .+ FROM auth_session WHERE account_id = \$1 AND revoked_at IS NULL ORDER BY rotated_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			s := newTestService(mock)

			sessions, err := s.GetUserSessions(context.Background(), tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, sessions, tt.wantCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
