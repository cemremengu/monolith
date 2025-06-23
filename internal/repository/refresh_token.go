package repository

import (
	"context"
	"errors"
	"time"

	"monolith/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrTokenNotFound = errors.New("refresh token not found")

type RefreshToken struct {
	ID        int64      `db:"id"`
	TokenHash string     `db:"token_hash"`
	AccountID uuid.UUID  `db:"account_id"`
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
	RevokedAt *time.Time `db:"revoked_at"`
}

type RefreshTokenRepository struct {
	db *database.DB
}

func NewRefreshTokenRepository(db *database.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(
	ctx context.Context,
	tokenHash string,
	accountID uuid.UUID,
	expiresAt time.Time,
) error {
	query := `
		INSERT INTO refresh_token (token_hash, account_id, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Pool.Exec(ctx, query, tokenHash, accountID, expiresAt)
	return err
}

func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	query := `
		SELECT id, token_hash, account_id, expires_at, created_at, revoked_at
		FROM refresh_token
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()
	`
	var token RefreshToken
	err := r.db.Pool.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.TokenHash,
		&token.AccountID,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE refresh_token
		SET revoked_at = NOW()
		WHERE token_hash = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, tokenHash)
	return err
}

func (r *RefreshTokenRepository) RevokeAllUserTokens(ctx context.Context, accountID uuid.UUID) error {
	query := `
		UPDATE refresh_token
		SET revoked_at = NOW()
		WHERE account_id = $1 AND revoked_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, accountID)
	return err
}

func (r *RefreshTokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	query := `
		DELETE FROM refresh_token
		WHERE expires_at < NOW() OR revoked_at < NOW() - INTERVAL '30 days'
	`
	_, err := r.db.Pool.Exec(ctx, query)
	return err
}
