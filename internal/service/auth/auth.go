package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"monolith/internal/config"
	"monolith/internal/database"
	"monolith/internal/database/dbsqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type Service struct {
	db             *database.DB
	queries        *dbsqlc.Queries
	securityConfig config.SecurityConfig
}

func NewService(db *database.DB, cfg config.SecurityConfig) *Service {
	return &Service{
		db:             db,
		queries:        db.Queries(),
		securityConfig: cfg,
	}
}

func (s *Service) CreateSession(ctx context.Context, req *CreateSessionRequest) (*Session, error) {
	token, hashedToken, err := createAndHashToken(s.securityConfig.SecretKey)
	if err != nil {
		return nil, err
	}

	row, err := s.queries.CreateSession(ctx, dbsqlc.CreateSessionParams{
		Token:     hashedToken,
		PrevToken: hashedToken,
		AccountID: req.AccountID,
		UserAgent: req.UserAgent,
		ClientIp:  req.ClientIP,
	})
	if err != nil {
		return nil, err
	}

	session := sessionFromAuthSession(row)
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

	row, err := s.queries.RotateSession(ctx, dbsqlc.RotateSessionParams{
		Token:     hashedToken,
		PrevToken: currentSession.Token,
		ID:        currentSession.ID,
	})
	if err != nil {
		return nil, err
	}

	session := sessionFromAuthSession(row)
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
	return s.queries.RevokeSession(ctx, dbsqlc.RevokeSessionParams{
		ID:        sessionID,
		AccountID: userID,
	})
}

func (s *Service) GetSessionByToken(ctx context.Context, unhashedToken string) (*Session, error) {
	hashedtoken := hashToken(unhashedToken, s.securityConfig.SecretKey)

	row, err := s.queries.GetSessionByToken(ctx, dbsqlc.GetSessionByTokenParams{
		Token:     hashedtoken,
		PrevToken: hashedtoken,
	})
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	session := Session{
		ID:        row.ID,
		Token:     row.Token,
		AccountID: row.AccountID,
		UserAgent: row.UserAgent,
		ClientIP:  row.ClientIp,
		CreatedAt: row.CreatedAt,
		RotatedAt: row.RotatedAt,
		RevokedAt: row.RevokedAt,
	}

	if session.RevokedAt != nil {
		return nil, ErrSessionRevoked
	}

	if session.CreatedAt.Before(s.createdAfterThreshold()) || session.RotatedAt.Before(s.rotatedAfterThreshold()) {
		return nil, ErrSessionExpired
	}

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

	row, err := s.queries.GetAuthContextByToken(ctx, dbsqlc.GetAuthContextByTokenParams{
		Token:     hashedtoken,
		PrevToken: hashedtoken,
		Status:    AccountStatusActive,
	})
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	authCtx := AuthContext{
		SessionID:      row.SessionID,
		SessionToken:   row.SessionToken,
		AccountID:      row.AccountID,
		AccountEmail:   row.AccountEmail,
		AccountIsAdmin: boolFromPgBool(row.AccountIsAdmin),
		AccountStatus:  row.AccountStatus,
		SessionCreated: row.SessionCreated,
		SessionRotated: row.SessionRotated,
		SessionRevoked: row.SessionRevoked,
	}

	if authCtx.SessionRevoked != nil {
		return nil, ErrSessionRevoked
	}

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
	return s.queries.RevokeAllUserSessions(ctx, accountID)
}

func (s *Service) GetSessionsByAccountID(ctx context.Context, accountID uuid.UUID) ([]Session, error) {
	rows, err := s.queries.GetSessionsByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, len(rows))
	for i, row := range rows {
		sessions[i] = Session{
			ID:        row.ID,
			Token:     row.Token,
			AccountID: row.AccountID,
			UserAgent: row.UserAgent,
			ClientIP:  row.ClientIp,
			CreatedAt: row.CreatedAt,
			RotatedAt: row.RotatedAt,
			RevokedAt: row.RevokedAt,
		}
	}
	return sessions, nil
}

func (s *Service) CleanupSessions(ctx context.Context) error {
	return s.queries.CleanupSessions(ctx)
}

func sessionFromAuthSession(row dbsqlc.AuthSession) Session {
	return Session{
		ID:        row.ID,
		Token:     row.Token,
		PrevToken: &row.PrevToken,
		AccountID: row.AccountID,
		UserAgent: row.UserAgent,
		ClientIP:  row.ClientIp,
		TokenSeen: row.TokenSeen,
		SeenAt:    row.SeenAt,
		CreatedAt: row.CreatedAt,
		RotatedAt: row.RotatedAt,
		RevokedAt: row.RevokedAt,
	}
}

func boolFromPgBool(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}
