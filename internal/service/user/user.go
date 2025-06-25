package user

import (
	"context"

	"monolith/internal/database"

	"github.com/georgysavva/scany/v2/pgxscan"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *Service) ValidatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) UserExists(ctx context.Context, email, username string) (bool, error) {
	var existingUser Account
	err := pgxscan.Get(ctx, s.db.Pool, &existingUser, `
		SELECT id FROM account WHERE email = $1 OR username = $2
	`, email, username)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) CreateAccount(ctx context.Context, req RegisterRequest) (*Account, error) {
	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var user Account
	err = pgxscan.Get(ctx, s.db.Pool, &user, `
		INSERT INTO account (username, email, name, password, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, NOW(), NOW()) 
		RETURNING id, username, email, name, is_admin, language, theme, timezone, 
		          last_seen_at, is_disabled, created_at, updated_at
	`, req.Username, req.Email, req.Name, hashedPassword)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) GetAccountByLogin(ctx context.Context, login string) (*Account, error) {
	var user Account
	err := pgxscan.Get(ctx, s.db.Pool, &user, `
		SELECT id, username, email, name, password, is_admin, language, theme, timezone, 
		       last_seen_at, is_disabled, created_at, updated_at
		FROM account 
		WHERE (email = $1 OR username = $1) AND is_disabled = FALSE
	`, login)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetAccountByID(ctx context.Context, userID string) (*Account, error) {
	var user Account
	err := pgxscan.Get(ctx, s.db.Pool, &user, `
		SELECT id, username, email, name, is_admin, language, theme, timezone, 
		       last_seen_at, is_disabled, created_at, updated_at
		FROM account 
		WHERE id = $1 AND is_disabled = FALSE
	`, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) UpdateLastSeen(ctx context.Context, userID string) error {
	_, err := s.db.Pool.Exec(ctx, `
		UPDATE account SET last_seen_at = NOW() WHERE id = $1
	`, userID)
	return err
}

func (s *Service) GetUsers(ctx context.Context) ([]Account, error) {
	var users []Account
	err := pgxscan.Select(ctx, s.db.Pool, &users, `
		SELECT id, username, email, name, is_admin, language, theme, timezone, 
		       last_seen_at, is_disabled, created_at, updated_at 
		FROM account 
		WHERE is_disabled = FALSE
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) GetUser(ctx context.Context, id string) (*Account, error) {
	var user Account
	err := pgxscan.Get(ctx, s.db.Pool, &user, `
		SELECT id, username, email, name, avatar, is_admin, language, theme, timezone, 
		       last_seen_at, is_disabled, created_at, updated_at 
		FROM account 
		WHERE id = $1 AND is_disabled = FALSE
	`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*Account, error) {
	var user Account
	err := pgxscan.Get(ctx, s.db.Pool, &user, `
		INSERT INTO account (username, name, email, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW()) 
		RETURNING id, username, email, name, is_admin, language, theme, timezone, 
		          last_seen_at, is_disabled, created_at, updated_at
	`, req.Username, req.Name, req.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) UpdateUser(ctx context.Context, id string, req CreateUserRequest) (*Account, error) {
	var user Account
	err := pgxscan.Get(ctx, s.db.Pool, &user, `
		UPDATE account 
		SET username = $1, name = $2, email = $3, updated_at = NOW() 
		WHERE id = $4 AND is_disabled = FALSE
		RETURNING id, username, email, name, is_admin, language, theme, timezone, 
		          last_seen_at, is_disabled, created_at, updated_at
	`, req.Username, req.Name, req.Email, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	_, err := s.db.Pool.Exec(ctx, `
		UPDATE account SET is_disabled = TRUE, updated_at = NOW() WHERE id = $1
	`, id)
	return err
}
