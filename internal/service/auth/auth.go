package auth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"monolith/internal/config"
	"monolith/internal/database"
	"monolith/internal/service/user"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const productionEnv = "production"

type Service struct {
	db             *database.DB
	tokenService   *TokenService
	sessionRepo    *SessionRepository
	securityConfig *config.SecurityConfig
	userService    *user.Service
}

func NewService(db *database.DB) *Service {
	return &Service{
		db:             db,
		tokenService:   NewTokenService(),
		sessionRepo:    NewSessionRepository(db),
		securityConfig: config.NewSecurityConfig(),
		userService:    user.NewService(db),
	}
}

func (s *Service) GenerateAndSetTokens(c echo.Context, userID, email string, isAdmin bool) error {
	accessToken, err := s.tokenService.GenerateAccessToken(userID, email, isAdmin)
	if err != nil {
		return err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return err
	}

	sessionID, err := s.sessionRepo.GenerateSessionID()
	if err != nil {
		return err
	}

	refreshTokenHash := HashToken(refreshToken, s.securityConfig.SecretKey)
	accountID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	userAgent := s.GetUserAgent(c)
	clientIP := s.GetClientIP(c)
	expiresAt := time.Now().Add(s.tokenService.RefreshTokenDuration())

	err = s.sessionRepo.CreateSession(
		c.Request().Context(),
		sessionID,
		refreshTokenHash,
		accountID,
		userAgent,
		clientIP,
		expiresAt,
	)
	if err != nil {
		return err
	}

	s.setCookies(c, accessToken, refreshToken, sessionID)
	return nil
}

func (s *Service) SetRefreshCookies(c echo.Context, accessToken, refreshToken string) {
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.tokenService.AccessTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(accessCookie)

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.tokenService.RefreshTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(refreshCookie)
}

func (s *Service) setCookies(c echo.Context, accessToken, refreshToken, sessionID string) {
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.tokenService.AccessTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(accessCookie)

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.tokenService.RefreshTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(refreshCookie)

	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.tokenService.RefreshTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(sessionCookie)
}

func (s *Service) Register(ctx context.Context, req user.RegisterRequest) (*user.Account, error) {
	if req.Password == "" || len(req.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	exists, err := s.userService.UserExists(ctx, req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	return s.userService.CreateAccount(ctx, req)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*user.Account, error) {
	user, err := s.userService.GetAccountByLogin(ctx, req.Login)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = s.userService.ValidatePassword(user.Password, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = s.userService.UpdateLastSeen(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) RefreshTokens(
	ctx context.Context,
	refreshToken,
	sessionID string,
) (*user.Account, string, string, error) {
	refreshTokenHash := HashToken(refreshToken, s.securityConfig.SecretKey)
	session, err := s.sessionRepo.GetSessionByToken(ctx, refreshTokenHash)

	if err != nil || session == nil || session.SessionID != sessionID {
		return nil, "", "", ErrSessionExpired
	}

	user, err := s.userService.GetAccountByID(ctx, session.AccountID.String())
	if err != nil {
		return nil, "", "", ErrUserNotFound
	}

	newAccessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return nil, "", "", err
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", err
	}

	newRefreshTokenHash := HashToken(newRefreshToken, s.securityConfig.SecretKey)
	newExpiresAt := time.Now().Add(s.tokenService.RefreshTokenDuration())
	err = s.sessionRepo.UpdateSessionToken(
		ctx,
		session.SessionID,
		newRefreshTokenHash,
		newExpiresAt,
	)
	if err != nil {
		return nil, "", "", err
	}

	return user, newAccessToken, newRefreshToken, nil
}

func (s *Service) ClearAuthCookies(c echo.Context) {
	cookies := []string{"access_token", "refresh_token", "session_id"}
	for _, cookieName := range cookies {
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    "",
			HttpOnly: true,
			Secure:   os.Getenv("ENV") == productionEnv,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
			Path:     "/",
		}
		c.SetCookie(cookie)
	}
}

func (s *Service) GetUserSessions(ctx context.Context, userID string) ([]SessionResponse, error) {
	accountID, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	sessions, err := s.sessionRepo.GetUserSessions(ctx, accountID)
	if err != nil {
		return nil, err
	}

	var response []SessionResponse
	for _, session := range sessions {
		response = append(response, SessionResponse{
			SessionID: session.SessionID,
			UserAgent: session.UserAgent,
			ClientIP:  session.ClientIP,
			CreatedAt: session.CreatedAt,
			RotatedAt: session.RotatedAt,
		})
	}

	return response, nil
}

func (s *Service) RevokeSession(ctx context.Context, userID, sessionID string) error {
	if userID != "" {
		accountID, err := uuid.Parse(userID)
		if err != nil {
			return ErrInvalidUserID
		}

		sessions, err := s.sessionRepo.GetUserSessions(ctx, accountID)
		if err != nil {
			return err
		}

		sessionFound := false
		for _, session := range sessions {
			if session.SessionID == sessionID {
				sessionFound = true
				break
			}
		}

		if !sessionFound {
			return ErrSessionNotFound
		}
	}

	return s.sessionRepo.RevokeSession(ctx, sessionID)
}

func (s *Service) RevokeAllOtherSessions(ctx context.Context, userID, currentSessionID string) (int, error) {
	accountID, err := uuid.Parse(userID)
	if err != nil {
		return 0, ErrInvalidUserID
	}

	sessions, err := s.sessionRepo.GetUserSessions(ctx, accountID)
	if err != nil {
		return 0, err
	}

	revokedCount := 0
	for _, session := range sessions {
		if session.SessionID != currentSessionID {
			err = s.sessionRepo.RevokeSession(ctx, session.SessionID)
			if err == nil {
				revokedCount++
			}
		}
	}

	return revokedCount, nil
}

// GetUserAgent extracts device information from User-Agent header.
func (s *Service) GetUserAgent(c echo.Context) string {
	userAgent := c.Request().Header.Get("User-Agent")
	if userAgent == "" {
		return "Unknown Device"
	}

	return userAgent
}

// GetClientIP extracts client IP address with support for proxy headers.
func (s *Service) GetClientIP(c echo.Context) string {
	if xff := c.Request().Header.Get("X-Forwarded-For"); xff != "" {
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if xri := c.Request().Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	return c.RealIP()
}
