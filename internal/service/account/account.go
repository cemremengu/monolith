package account

import (
	"context"

	"monolith/internal/database"
	"monolith/internal/util"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
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

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var account Account
	err = pgxscan.Get(ctx, s.db.Pool, &account, `
		INSERT INTO account (username, email, name, password, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, NOW(), NOW()) 
		RETURNING id, username, email, name, is_admin, language, theme, timezone, 
		          last_seen_at, is_disabled, created_at, updated_at
	`, req.Username, req.Email, req.Name, hashedPassword)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *Service) GetAccountByLogin(ctx context.Context, login string) (*Account, error) {
	var account Account
	err := pgxscan.Get(ctx, s.db.Pool, &account, `
		SELECT id, username, email, name, password, is_admin, language, theme, timezone, 
		       last_seen_at, is_disabled, created_at, updated_at
		FROM account 
		WHERE (email = $1 OR username = $1) AND is_disabled = FALSE
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
		       last_seen_at, is_disabled, created_at, updated_at
		FROM account 
		WHERE id = $1 AND is_disabled = FALSE
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
		WHERE id = $4 AND is_disabled = FALSE
		RETURNING id, username, email, name, is_admin, language, theme, timezone, 
		          last_seen_at, is_disabled, created_at, updated_at
	`, req.Language, req.Theme, req.Timezone, accountID)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
