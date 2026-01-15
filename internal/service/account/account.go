package account

import (
	"context"
	"fmt"
	"strings"

	"monolith/internal/database"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) ValidatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) UserExists(ctx context.Context, email, username string) (bool, error) {
	var existingAccount Account
	err := pgxscan.Get(ctx, s.db.Pool, &existingAccount, `
		SELECT id FROM account WHERE email = $1 OR username = $2
	`, email, username)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*Account, error) {
	if req.Password == "" || len(req.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	exists, err := s.UserExists(ctx, req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var account Account
	err = pgxscan.Get(ctx, tx, &account, `
		INSERT INTO account (username, email, name, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, username, email, name, is_admin, language, theme, timezone,
		          last_seen_at, status, created_at, updated_at
	`, req.Username, req.Email, req.Name, hashedPassword)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *Service) GetAccountByLogin(ctx context.Context, login string) (*Account, error) {
	var account Account
	err := pgxscan.Get(ctx, s.db.Pool, &account, `
		SELECT id, username, email, name, password, is_admin, language, theme, timezone,
		       last_seen_at, status, created_at, updated_at
		FROM account
		WHERE (email = $1 OR username = $1) AND status = 'active'
	`, login)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *Service) GetAccountByID(ctx context.Context, accountID uuid.UUID) (*Account, error) {
	var account Account
	err := pgxscan.Get(ctx, s.db.Pool, &account, `
		SELECT id, username, email, name, is_admin, language, theme, timezone,
		       last_seen_at, status, created_at, updated_at
		FROM account
		WHERE id = $1 AND status = 'active'
	`, accountID)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *Service) UpdateLastSeen(ctx context.Context, accountID uuid.UUID) error {
	_, err := s.db.Pool.Exec(ctx, `
		UPDATE account SET last_seen_at = NOW() WHERE id = $1
	`, accountID)
	return err
}

func (s *Service) UpdatePreferences(
	ctx context.Context,
	accountID uuid.UUID,
	req UpdatePreferencesRequest,
) (*Account, error) {
	var account Account
	err := pgxscan.Get(ctx, s.db.Pool, &account, `
		UPDATE account
		SET language = $1, theme = $2, timezone = $3, updated_at = NOW()
		WHERE id = $4 AND status = 'active'
		RETURNING id, username, email, name, is_admin, language, theme, timezone,
		          last_seen_at, status, created_at, updated_at
	`, req.Language, req.Theme, req.Timezone, accountID)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *Service) GetAccounts(ctx context.Context) ([]Account, error) {
	var accounts []Account
	err := pgxscan.Select(ctx, s.db.Pool, &accounts, `
		SELECT id, username, email, name, avatar, is_admin, language, theme, timezone,
		       last_seen_at, status, created_at, updated_at
		FROM account
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *Service) GetAccount(ctx context.Context, id uuid.UUID) (*Account, error) {
	var account Account
	err := pgxscan.Get(ctx, s.db.Pool, &account, `
		SELECT id, username, email, name, avatar, is_admin, language, theme, timezone,
		       last_seen_at, status, created_at, updated_at
		FROM account
		WHERE id = $1 AND status = 'active'
	`, id)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *Service) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	var account Account
	var hashedPassword *string
	status := "pending"
	isAdmin := false

	if req.Password != nil && *req.Password != "" {
		hashed, err := hashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		hashedPassword = &hashed
		status = "active"
	}

	if req.IsAdmin != nil {
		isAdmin = *req.IsAdmin
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	err = pgxscan.Get(ctx, tx, &account, `
		INSERT INTO account (username, name, email, password, is_admin, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, username, email, name, avatar, is_admin, language, theme, timezone,
		          last_seen_at, status, created_at, updated_at
	`, req.Username, req.Name, req.Email, hashedPassword, isAdmin, status)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *Service) InviteUsers(ctx context.Context, req InviteUsersRequest) (*InviteUsersResponse, error) {
	response := &InviteUsersResponse{
		Success: []Account{},
		Failed:  []InviteUserFailure{},
	}

	for _, email := range req.Emails {
		username := deriveUsernameFromEmail(email)

		var existingAccount Account
		err := pgxscan.Get(ctx, s.db.Pool, &existingAccount, `
			SELECT id FROM account WHERE email = $1 OR username = $2
		`, email, username)

		if err == nil {
			response.Failed = append(response.Failed, InviteUserFailure{
				Email:  email,
				Reason: "User already exists",
			})
			continue
		}

		var account Account
		err = pgxscan.Get(ctx, s.db.Pool, &account, `
			INSERT INTO account (username, email, name, is_admin, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, 'pending', NOW(), NOW())
			RETURNING id, username, email, name, avatar, is_admin, language, theme, timezone,
			          last_seen_at, status, created_at, updated_at
		`, username, email, username, req.IsAdmin)
		if err != nil {
			response.Failed = append(response.Failed, InviteUserFailure{
				Email:  email,
				Reason: fmt.Sprintf("Failed to create user: %v", err),
			})
			continue
		}

		response.Success = append(response.Success, account)
	}

	return response, nil
}

func (s *Service) UpdateAccount(ctx context.Context, id uuid.UUID, req UpdateAccountRequest) (*Account, error) {
	var account Account
	err := pgxscan.Get(ctx, s.db.Pool, &account, `
		UPDATE account
		SET username = $1, name = $2, email = $3, updated_at = NOW()
		WHERE id = $4 AND status = 'active'
		RETURNING id, username, email, name, avatar, is_admin, language, theme, timezone,
		          last_seen_at, status, created_at, updated_at
	`, req.Username, req.Name, req.Email, id)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *Service) DisableAccount(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Pool.Exec(ctx, `
		UPDATE account SET status = 'disabled', updated_at = NOW() WHERE id = $1
	`, id)
	return err
}

func (s *Service) EnableAccount(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Pool.Exec(ctx, `
		UPDATE account SET status = 'active', updated_at = NOW() WHERE id = $1
	`, id)
	return err
}

// TODO: this should trigger deleting other stuff
func (s *Service) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Pool.Exec(ctx, `
		DELETE FROM account WHERE id = $1
	`, id)
	return err
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func deriveUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		username := strings.ToLower(parts[0])
		username = strings.ReplaceAll(username, ".", "")
		username = strings.ReplaceAll(username, "+", "")
		return username
	}
	return email
}

func (s *Service) ChangePassword(ctx context.Context, accountID uuid.UUID, req ChangePasswordRequest) error {
	if req.NewPassword == "" || len(req.NewPassword) < 8 {
		return ErrPasswordTooShort
	}

	var account Account
	err := pgxscan.Get(ctx, s.db.Pool, &account, `
		SELECT id, password
		FROM account
		WHERE id = $1 AND status = 'active'
	`, accountID)
	if err != nil {
		return err
	}

	if err := s.ValidatePassword(account.Password, req.CurrentPassword); err != nil {
		return ErrInvalidPassword
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	_, err = s.db.Pool.Exec(ctx, `
		UPDATE account SET password = $1, updated_at = NOW() WHERE id = $2
	`, hashedPassword, accountID)
	return err
}
