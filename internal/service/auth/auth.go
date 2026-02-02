package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"monolith/internal/config"
	"monolith/internal/database"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Service struct {
	db             *database.DB
	securityConfig config.SecurityConfig
}

func NewService(db *database.DB, cfg config.SecurityConfig) *Service {
	return &Service{
		db:             db,
		securityConfig: cfg,
	}
}

func (s *Service) CreateSession(ctx context.Context, req *CreateSessionRequest) (*Session, error) {
	token, hashedToken, err := createAndHashToken(s.securityConfig.SecretKey)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO auth_session (token, prev_token, account_id, user_agent, client_ip)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`

	var session Session
	err = pgxscan.Get(ctx, s.db.Pool, &session, query, hashedToken, hashedToken, req.AccountID, req.UserAgent, req.ClientIP)
	if err != nil {
		return nil, err
	}

	session.UnhashedToken = token

	return &session, nil
}

func (s *Service) SetSessionCookies(c echo.Context, session *Session) {
	c.SetCookie(&http.Cookie{
		Name:     s.securityConfig.LoginCookieName,
		Value:    session.UnhashedToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.securityConfig.LoginMaximumLifetimeDuration.Seconds()),
		Path:     "/",
	})

	expiry := session.NextRotation(time.Duration(s.securityConfig.TokenRotationIntervalMinutes) * time.Minute)

	c.SetCookie(&http.Cookie{
		Name:     "session_expiry",
		Value:    strconv.FormatInt(expiry.Unix(), 10),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.securityConfig.LoginMaximumLifetimeDuration.Seconds()),
		Path:     "/",
	})
}

func (s *Service) RotateSession(ctx context.Context, req *RotateSessionRequest) (*Session, error) {
	currentSession, err := s.GetSessionByToken(ctx, req.UnhashedToken)
	if err != nil {
		return nil, err
	}

	newToken, hashedToken, err := createAndHashToken(s.securityConfig.SecretKey)
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE auth_session
		SET token = $1, prev_token = $2, rotated_at = NOW(), token_seen = FALSE, seen_at = NULL
		WHERE id = $3
		RETURNING *
	`

	var session Session
	err = pgxscan.Get(ctx, s.db.Pool, &session, query, hashedToken, currentSession.Token, currentSession.ID)
	if err != nil {
		return nil, err
	}

	session.UnhashedToken = newToken

	return &session, nil
}

func (s *Service) ClearAuthCookies(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     s.securityConfig.LoginCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Path:     "/",
	})

	c.SetCookie(&http.Cookie{
		Name:     "session_expiry",
		Value:    "",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Path:     "/",
	})
}

func (s *Service) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]UserSession, error) {
	sessions, err := s.GetSessionsByAccountID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var response []UserSession
	for _, session := range sessions {
		response = append(response, UserSession{
			ID:        session.ID,
			UserAgent: session.UserAgent,
			ClientIP:  session.ClientIP,
			CreatedAt: session.CreatedAt,
			RotatedAt: session.RotatedAt,
		})
	}

	return response, nil
}

func (s *Service) RevokeSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	query := `
		UPDATE auth_session
		SET revoked_at = NOW()
		WHERE id = $1 and account_id = $2
	`
	_, err := s.db.Pool.Exec(ctx, query, sessionID, userID)
	return err
}

func (s *Service) GetSessionByToken(ctx context.Context, unhashedToken string) (*Session, error) {
	hashedtoken := hashToken(unhashedToken, s.securityConfig.SecretKey)

	query := `
		SELECT id, token, account_id, user_agent, client_ip, created_at, rotated_at, revoked_at
		FROM auth_session
		WHERE token = $1 OR prev_token = $2
	`
	var session Session
	err := pgxscan.Get(ctx, s.db.Pool, &session, query, hashedtoken, hashedtoken)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	if session.RevokedAt != nil {
		return nil, ErrSessionRevoked
	}

	if session.CreatedAt.Before(s.createdAfterThreshold()) || session.RotatedAt.Before(s.rotatedAfterThreshold()) {
		return nil, ErrSessionExpired
	}

	// if session.NeedsRotation(time.Duration(s.securityConfig.TokenRotationIntervalMinutes) * time.Minute) {
	// 	return nil, ErrSessionNeedsRotation
	// }

	return &session, nil
}

// GetAuthContextByToken retrieves all authentication context (session, account, workspace) in a single query.
// It performs the following security validations:
// - Session token validation (matches current or previous token)
// - Account status verification (must be active)
// - Session revocation check
// - Session expiration check (both created and rotated timestamps)
//
// Note: The account status check is performed in the JOIN condition rather than a separate WHERE clause
// for security purposes - this prevents information leakage about account states by returning a generic
// "session not found" error for both non-existent sessions and inactive accounts.
//
// Returns:
// - ErrSessionNotFound if session doesn't exist or account is not active
// - ErrSessionRevoked if session has been revoked
// - ErrSessionExpired if session has expired
func (s *Service) GetAuthContextByToken(ctx context.Context, unhashedToken string) (*AuthContext, error) {
	hashedtoken := hashToken(unhashedToken, s.securityConfig.SecretKey)

	query := `
		SELECT 
			s.id as session_id,
			s.token as session_token,
			s.account_id,
			a.email as account_email,
			a.is_admin as account_is_admin,
			a.status as account_status,
			s.created_at as session_created,
			s.rotated_at as session_rotated,
			s.revoked_at as session_revoked
		FROM auth_session s
		INNER JOIN account a ON s.account_id = a.id AND a.status = $3
		WHERE (s.token = $1 OR s.prev_token = $2)
	`

	var authCtx AuthContext
	err := pgxscan.Get(ctx, s.db.Pool, &authCtx, query, hashedtoken, hashedtoken, AccountStatusActive)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			// Return generic error to avoid leaking account state information
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	// Check if session is revoked
	if authCtx.SessionRevoked != nil {
		return nil, ErrSessionRevoked
	}

	// Check if session is expired
	if authCtx.SessionCreated.Before(s.createdAfterThreshold()) || authCtx.SessionRotated.Before(s.rotatedAfterThreshold()) {
		return nil, ErrSessionExpired
	}

	return &authCtx, nil
}

func (s *Service) createdAfterThreshold() time.Time {
	return time.Now().Add(-s.securityConfig.LoginMaximumLifetimeDuration)
}

func (s *Service) rotatedAfterThreshold() time.Time {
	return time.Now().Add(-s.securityConfig.LoginMaximumInactiveLifetimeDuration)
}

func (s *Service) RevokeAllUserSessions(ctx context.Context, accountID uuid.UUID) error {
	query := `
		UPDATE auth_session
		SET revoked_at = NOW()
		WHERE account_id = $1 AND revoked_at IS NULL
	`
	_, err := s.db.Pool.Exec(ctx, query, accountID)
	return err
}

func (s *Service) GetSessionsByAccountID(ctx context.Context, accountID uuid.UUID) ([]Session, error) {
	query := `
		SELECT id, token, account_id, user_agent, client_ip, created_at, rotated_at, revoked_at
		FROM auth_session
		WHERE account_id = $1 AND revoked_at IS NULL
		ORDER BY rotated_at DESC
	`
	var sessions []Session
	err := pgxscan.Select(ctx, s.db.Pool, &sessions, query, accountID)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (s *Service) CleanupSessions(ctx context.Context) error {
	query := `
		DELETE FROM auth_session
		WHERE revoked_at IS NOT NULL
	`
	_, err := s.db.Pool.Exec(ctx, query)
	return err
}
