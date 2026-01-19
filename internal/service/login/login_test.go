package login

import (
	"context"
	"errors"
	"testing"
	"time"

	"monolith/internal/database"
	"monolith/internal/service/account"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func stringPtr(s string) *string {
	return &s
}

func TestService_Login(t *testing.T) {
	accountID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	now := time.Now()

	tests := []struct {
		name      string
		req       UserLoginRequest
		setupMock func(mock pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name: "successful login",
			req: UserLoginRequest{
				Login:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "username", "email", "name", "password", "is_admin", "language", "theme", "timezone",
					"last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(
					accountID, "testuser", "test@example.com", stringPtr("Test User"),
					string(hashedPassword), false, stringPtr("en"), stringPtr("light"),
					stringPtr("UTC"), nil, "active", now, now,
				)
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("test@example.com").
					WillReturnRows(rows)

				mock.ExpectExec(`UPDATE account SET last_seen_at = NOW\(\) WHERE id = \$1`).
					WithArgs(accountID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			req: UserLoginRequest{
				Login:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("nonexistent@example.com").
					WillReturnError(database.ErrNoRows)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			req: UserLoginRequest{
				Login:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "username", "email", "name", "password", "is_admin", "language", "theme", "timezone",
					"last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(
					accountID, "testuser", "test@example.com", stringPtr("Test User"),
					string(hashedPassword), false, stringPtr("en"), stringPtr("light"),
					stringPtr("UTC"), nil, "active", now, now,
				)
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "login by username",
			req: UserLoginRequest{
				Login:    "testuser",
				Password: "password123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "username", "email", "name", "password", "is_admin", "language", "theme", "timezone",
					"last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(
					accountID, "testuser", "test@example.com", stringPtr("Test User"),
					string(hashedPassword), false, stringPtr("en"), stringPtr("light"),
					stringPtr("UTC"), nil, "active", now, now,
				)
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("testuser").
					WillReturnRows(rows)

				mock.ExpectExec(`UPDATE account SET last_seen_at = NOW\(\) WHERE id = \$1`).
					WithArgs(accountID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name: "update last seen fails",
			req: UserLoginRequest{
				Login:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "username", "email", "name", "password", "is_admin", "language", "theme", "timezone",
					"last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(
					accountID, "testuser", "test@example.com", stringPtr("Test User"),
					string(hashedPassword), false, stringPtr("en"), stringPtr("light"),
					stringPtr("UTC"), nil, "active", now, now,
				)
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("test@example.com").
					WillReturnRows(rows)

				mock.ExpectExec(`UPDATE account SET last_seen_at = NOW\(\) WHERE id = \$1`).
					WithArgs(accountID).
					WillReturnError(errors.New("database error"))
			},
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			accountSvc := account.NewService(db)
			s := NewService(db, accountSvc)

			acc, err := s.Login(context.Background(), tt.req)
			if tt.wantErr != nil {
				require.Error(t, err)
				if errors.Is(tt.wantErr, ErrInvalidCredentials) {
					assert.ErrorIs(t, err, ErrInvalidCredentials)
				}
				assert.Nil(t, acc)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, acc)
				assert.Equal(t, accountID, acc.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
