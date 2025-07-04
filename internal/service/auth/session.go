package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"monolith/internal/database"

	"github.com/google/uuid"
)

const sessionIDLength = 32

type Session struct {
	ID         int64      `db:"id"`
	SessionID  string     `db:"session_id"`
	TokenHash  string     `db:"token_hash"`
	AccountID  uuid.UUID  `db:"account_id"`
	DeviceInfo string     `db:"device_info"`
	IPAddress  string     `db:"ip_address"`
	ExpiresAt  time.Time  `db:"expires_at"`
	CreatedAt  time.Time  `db:"created_at"`
	RotatedAt  time.Time  `db:"rotated_at"`
	RevokedAt  *time.Time `db:"revoked_at"`
}

type SessionRepository struct {
	db *database.DB
}

func NewSessionRepository(db *database.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) GenerateSessionID() (string, error) {
	bytes := make([]byte, sessionIDLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (r *SessionRepository) CreateSession(
	ctx context.Context,
	sessionID, tokenHash string,
	accountID uuid.UUID,
	deviceInfo, ipAddress string,
	expiresAt time.Time,
) error {
	query := `
		INSERT INTO session (session_id, token_hash, account_id, device_info, ip_address, expires_at, rotated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	_, err := r.db.Pool.Exec(ctx, query, sessionID, tokenHash, accountID, deviceInfo, ipAddress, expiresAt)
	return err
}

func (r *SessionRepository) GetSessionByToken(ctx context.Context, tokenHash string) (*Session, error) {
	query := `
		SELECT id, session_id, token_hash, account_id, device_info, ip_address, expires_at, created_at, rotated_at, revoked_at
		FROM session
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()
	`
	var session Session
	err := r.db.Pool.QueryRow(ctx, query, tokenHash).Scan(
		&session.ID,
		&session.SessionID,
		&session.TokenHash,
		&session.AccountID,
		&session.DeviceInfo,
		&session.IPAddress,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.RotatedAt,
		&session.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetSessionByTokenWithTimeout(
	ctx context.Context,
	tokenHash string,
	sessionTimeout time.Duration,
) (*Session, error) {
	var query string
	var args []any

	if sessionTimeout > 0 {
		query = `
			SELECT id, session_id, token_hash, account_id, device_info, ip_address, expires_at, created_at, rotated_at, revoked_at
			FROM session
			WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW() AND created_at > $2
		`
		timeoutThreshold := time.Now().Add(-sessionTimeout)
		args = []any{tokenHash, timeoutThreshold}
	} else {
		query = `
			SELECT id, session_id, token_hash, account_id, device_info, ip_address, expires_at, created_at, rotated_at, revoked_at
			FROM session
			WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()
		`
		args = []any{tokenHash}
	}

	var session Session
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&session.ID,
		&session.SessionID,
		&session.TokenHash,
		&session.AccountID,
		&session.DeviceInfo,
		&session.IPAddress,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.RotatedAt,
		&session.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) UpdateSessionToken(
	ctx context.Context,
	sessionID, newTokenHash string,
	expiresAt time.Time,
) error {
	query := `
		UPDATE session
		SET token_hash = $1, expires_at = $2, rotated_at = NOW()
		WHERE session_id = $3 AND revoked_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, newTokenHash, expiresAt, sessionID)
	return err
}

func (r *SessionRepository) RevokeSession(ctx context.Context, sessionID string) error {
	query := `
		UPDATE session
		SET revoked_at = NOW()
		WHERE session_id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, sessionID)
	return err
}

func (r *SessionRepository) RevokeAllUserSessions(ctx context.Context, accountID uuid.UUID) error {
	query := `
		UPDATE session
		SET revoked_at = NOW()
		WHERE account_id = $1 AND revoked_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, accountID)
	return err
}

func (r *SessionRepository) GetUserSessions(ctx context.Context, accountID uuid.UUID) ([]Session, error) {
	query := `
		SELECT id, session_id, token_hash, account_id, device_info, ip_address, expires_at, created_at, rotated_at, revoked_at
		FROM session
		WHERE account_id = $1 AND revoked_at IS NULL AND expires_at > NOW()
		ORDER BY rotated_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var session Session
		scanErr := rows.Scan(
			&session.ID,
			&session.SessionID,
			&session.TokenHash,
			&session.AccountID,
			&session.DeviceInfo,
			&session.IPAddress,
			&session.ExpiresAt,
			&session.CreatedAt,
			&session.RotatedAt,
			&session.RevokedAt,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		sessions = append(sessions, session)
	}
	return sessions, rows.Err()
}

func (r *SessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	query := `
		DELETE FROM session
		WHERE expires_at < NOW() OR revoked_at < NOW() - INTERVAL '30 days'
	`
	_, err := r.db.Pool.Exec(ctx, query)
	return err
}
