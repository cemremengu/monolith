package user

import (
	"context"

	"monolith/internal/auth"
	"monolith/internal/database"

	"github.com/georgysavva/scany/v2/pgxscan"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Name     string `json:"name"     validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
}

func (s *Service) GetUsers(ctx context.Context) ([]auth.UserAccount, error) {
	var users []auth.UserAccount
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

func (s *Service) GetUser(ctx context.Context, id string) (*auth.UserAccount, error) {
	var user auth.UserAccount
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

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*auth.UserAccount, error) {
	var user auth.UserAccount
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

func (s *Service) UpdateUser(ctx context.Context, id string, req CreateUserRequest) (*auth.UserAccount, error) {
	var user auth.UserAccount
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
