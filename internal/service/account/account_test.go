package account

import (
	"context"
	"testing"
	"time"

	"monolith/internal/database"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func stringPtr(s string) *string {
	return &s
}

func TestService_ValidatePassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
	}{
		{
			name:           "valid password",
			hashedPassword: string(hashedPassword),
			password:       "password123",
			wantErr:        false,
		},
		{
			name:           "wrong password",
			hashedPassword: string(hashedPassword),
			password:       "wrongpassword",
			wantErr:        true,
		},
		{
			name:           "empty password",
			hashedPassword: string(hashedPassword),
			password:       "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			db := &database.DB{Pool: mock}
			s := NewService(db)

			err := s.ValidatePassword(tt.hashedPassword, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_UserExists(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		username  string
		setupMock func(mock pgxmock.PgxPoolIface)
		want      bool
		wantErr   bool
	}{
		{
			name:     "user exists",
			email:    "test@example.com",
			username: "testuser",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(uuid.New())
				mock.ExpectQuery(`SELECT id FROM account WHERE email = \$1 OR username = \$2`).
					WithArgs("test@example.com", "testuser").
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:     "user not found",
			email:    "notfound@example.com",
			username: "notfound",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT id FROM account WHERE email = \$1 OR username = \$2`).
					WithArgs("notfound@example.com", "notfound").
					WillReturnError(database.ErrNoRows)
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			s := NewService(db)

			got, err := s.UserExists(context.Background(), tt.email, tt.username)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_Register(t *testing.T) {
	tests := []struct {
		name      string
		req       RegisterRequest
		setupMock func(mock pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name: "password too short",
			req: RegisterRequest{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "short",
				Name:     "New User",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {},
			wantErr:   ErrPasswordTooShort,
		},
		{
			name: "empty password",
			req: RegisterRequest{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "",
				Name:     "New User",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {},
			wantErr:   ErrPasswordTooShort,
		},
		{
			name: "user already exists",
			req: RegisterRequest{
				Username: "existinguser",
				Email:    "existing@example.com",
				Password: "password123",
				Name:     "Existing User",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(uuid.New())
				mock.ExpectQuery(`SELECT id FROM account WHERE email = \$1 OR username = \$2`).
					WithArgs("existing@example.com", "existinguser").
					WillReturnRows(rows)
			},
			wantErr: ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			s := NewService(db)

			account, err := s.Register(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, account)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tt.req.Username, account.Username)
				assert.Equal(t, tt.req.Email, account.Email)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_GetAccountByLogin(t *testing.T) {
	accountID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		login     string
		setupMock func(mock pgxmock.PgxPoolIface)
		want      *Account
		wantErr   bool
	}{
		{
			name:  "found by email",
			login: "test@example.com",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "username", "email", "name", "password", "is_admin", "language", "theme", "timezone",
					"last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(
					accountID, "testuser", "test@example.com", stringPtr("Test User"),
					"hashedpassword", false, stringPtr("en"), stringPtr("light"),
					stringPtr("UTC"), nil, "active", now, now,
				)
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want: &Account{
				ID:       accountID,
				Username: "testuser",
				Email:    "test@example.com",
				Status:   "active",
			},
			wantErr: false,
		},
		{
			name:  "found by username",
			login: "testuser",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "username", "email", "name", "password", "is_admin", "language", "theme", "timezone",
					"last_seen_at", "status", "created_at", "updated_at",
				}).AddRow(
					accountID, "testuser", "test@example.com", stringPtr("Test User"),
					"hashedpassword", false, stringPtr("en"), stringPtr("light"),
					stringPtr("UTC"), nil, "active", now, now,
				)
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			want: &Account{
				ID:       accountID,
				Username: "testuser",
				Email:    "test@example.com",
				Status:   "active",
			},
			wantErr: false,
		},
		{
			name:  "not found",
			login: "nonexistent",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT .+ FROM account WHERE \(email = \$1 OR username = \$1\) AND status = 'active'`).
					WithArgs("nonexistent").
					WillReturnError(database.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			s := NewService(db)

			got, err := s.GetAccountByLogin(context.Background(), tt.login)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_GetAccountByID(t *testing.T) {
	accountID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		accountID uuid.UUID
		setupMock func(mock pgxmock.PgxPoolIface)
		wantErr   bool
	}{
		{
			name:      "account found",
			accountID: accountID,
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
			wantErr: false,
		},
		{
			name:      "account not found",
			accountID: uuid.New(),
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT .+ FROM account WHERE id = \$1 AND status = 'active'`).
					WithArgs(pgxmock.AnyArg()).
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

			db := &database.DB{Pool: mock}
			s := NewService(db)

			got, err := s.GetAccountByID(context.Background(), tt.accountID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.accountID, got.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_ChangePassword(t *testing.T) {
	accountID := uuid.New()
	hashedCurrentPassword, _ := bcrypt.GenerateFromPassword([]byte("currentpass123"), 12)
	now := time.Now()

	tests := []struct {
		name      string
		accountID uuid.UUID
		req       ChangePasswordRequest
		setupMock func(mock pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:      "successful password change",
			accountID: accountID,
			req: ChangePasswordRequest{
				CurrentPassword: "currentpass123",
				NewPassword:     "newpassword123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "password"}).
					AddRow(accountID, string(hashedCurrentPassword))
				mock.ExpectQuery(`SELECT id, password FROM account WHERE id = \$1 AND status = 'active'`).
					WithArgs(accountID).
					WillReturnRows(rows)

				mock.ExpectExec(`UPDATE account SET password = \$1, updated_at = NOW\(\) WHERE id = \$2`).
					WithArgs(pgxmock.AnyArg(), accountID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name:      "new password too short",
			accountID: accountID,
			req: ChangePasswordRequest{
				CurrentPassword: "currentpass123",
				NewPassword:     "short",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {},
			wantErr:   ErrPasswordTooShort,
		},
		{
			name:      "wrong current password",
			accountID: accountID,
			req: ChangePasswordRequest{
				CurrentPassword: "wrongpassword",
				NewPassword:     "newpassword123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "password"}).
					AddRow(accountID, string(hashedCurrentPassword))
				mock.ExpectQuery(`SELECT id, password FROM account WHERE id = \$1 AND status = 'active'`).
					WithArgs(accountID).
					WillReturnRows(rows)
			},
			wantErr: ErrInvalidPassword,
		},
		{
			name:      "account not found",
			accountID: uuid.New(),
			req: ChangePasswordRequest{
				CurrentPassword: "currentpass123",
				NewPassword:     "newpassword123",
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT id, password FROM account WHERE id = \$1 AND status = 'active'`).
					WithArgs(pgxmock.AnyArg()).
					WillReturnError(database.ErrNoRows)
			},
			wantErr: database.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			_ = now
			tt.setupMock(mock)

			db := &database.DB{Pool: mock}
			s := NewService(db)

			err := s.ChangePassword(context.Background(), tt.accountID, tt.req)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "expected %v, got %v", tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_UpdatePreferences(t *testing.T) {
	accountID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		accountID uuid.UUID
		req       UpdatePreferencesRequest
		setupMock func(mock pgxmock.PgxPoolIface)
		wantErr   bool
	}{
		{
			name:      "successful update",
			accountID: accountID,
			req: UpdatePreferencesRequest{
				Language: stringPtr("fr"),
				Theme:    stringPtr("dark"),
				Timezone: stringPtr("Europe/Paris"),
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
			wantErr: false,
		},
		{
			name:      "account not found",
			accountID: uuid.New(),
			req: UpdatePreferencesRequest{
				Language: stringPtr("fr"),
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`UPDATE account SET language = \$1, theme = \$2, timezone = \$3, updated_at = NOW\(\) WHERE id = \$4 AND status = 'active' RETURNING`).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
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

			db := &database.DB{Pool: mock}
			s := NewService(db)

			got, err := s.UpdatePreferences(context.Background(), tt.accountID, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
