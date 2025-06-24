package user

import (
	"context"
	"time"

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

type UserAccount struct {
	ID         string     `json:"id"               db:"id"`
	Username   string     `json:"username"         db:"username"`
	Email      string     `json:"email"            db:"email"`
	Name       *string    `json:"name"             db:"name"`
	Password   string     `json:"-"                db:"password"`
	Avatar     *string    `json:"avatar,omitempty" db:"avatar"`
	IsAdmin    bool       `json:"isAdmin"          db:"is_admin"`
	Language   *string    `json:"language"         db:"language"`
	Theme      *string    `json:"theme"            db:"theme"`
	Timezone   *string    `json:"timezone"         db:"timezone"`
	LastSeenAt *time.Time `json:"lastSeenAt"       db:"last_seen_at"`
	IsDisabled bool       `json:"isDisabled"       db:"is_disabled"`
	CreatedAt  time.Time  `json:"createdAt"        db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt"        db:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Name     string `json:"name"     validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
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
	var existingUser UserAccount
	err := pgxscan.Get(ctx, s.db.Pool, &existingUser, `
		SELECT id FROM account WHERE email = $1 OR username = $2
	`, email, username)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *Service) CreateAccount(ctx context.Context, req RegisterRequest) (*UserAccount, error) {
	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var user UserAccount
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

func (s *Service) GetAccountByLogin(ctx context.Context, login string) (*UserAccount, error) {
	var user UserAccount
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

func (s *Service) GetAccountByID(ctx context.Context, userID string) (*UserAccount, error) {
	var user UserAccount
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

func (s *Service) GetUsers(ctx context.Context) ([]UserAccount, error) {
	var users []UserAccount
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

func (s *Service) GetUser(ctx context.Context, id string) (*UserAccount, error) {
	var user UserAccount
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

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserAccount, error) {
	var user UserAccount
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

func (s *Service) UpdateUser(ctx context.Context, id string, req CreateUserRequest) (*UserAccount, error) {
	var user UserAccount
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
