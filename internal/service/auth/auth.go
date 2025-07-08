package auth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"monolith/internal/config"
	"monolith/internal/database"
	"monolith/internal/service/account"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const productionEnv = "production"

type Service struct {
	db             *database.DB
	tokenService   *TokenService
	sessionRepo    *SessionRepository
	securityConfig *config.SecurityConfig
	accountService *account.Service
}

func NewService(db *database.DB) *Service {
	return &Service{
		db:             db,
		tokenService:   NewTokenService(),
		sessionRepo:    NewSessionRepository(db),
		securityConfig: config.NewSecurityConfig(),
		accountService: account.NewService(db),
	}
}

func (s *Service) GenerateAndSetTokens(c echo.Context, userID uuid.UUID, email string, isAdmin bool) error {
	accessToken, err := s.tokenService.GenerateAccessToken(userID, email, isAdmin)
	if err != nil {
		return err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return err
	}

	sessionID := uuid.New()

	refreshTokenHash := HashToken(refreshToken, s.securityConfig.SecretKey)

	userAgent := s.GetUserAgent(c)
	clientIP := s.GetClientIP(c)
	expiresAt := time.Now().Add(s.tokenService.RefreshTokenDuration())

	err = s.sessionRepo.CreateSession(
		c.Request().Context(),
		sessionID,
		refreshTokenHash,
		userID,
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

func (s *Service) setCookies(c echo.Context, accessToken string, refreshToken string, sessionID uuid.UUID) {
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
		Value:    sessionID.String(),
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.tokenService.RefreshTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(sessionCookie)
}

func (s *Service) Register(ctx context.Context, req account.RegisterRequest) (*account.Account, error) {
	if req.Password == "" || len(req.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	exists, err := s.accountService.UserExists(ctx, req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	return s.accountService.Register(ctx, req)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*account.Account, error) {
	account, err := s.accountService.GetAccountByLogin(ctx, req.Login)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = s.accountService.ValidatePassword(account.Password, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = s.accountService.UpdateLastSeen(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *Service) RefreshTokens(
	ctx context.Context,
	refreshToken string,
	sessionID uuid.UUID,
) (*account.Account, string, string, error) {
	refreshTokenHash := HashToken(refreshToken, s.securityConfig.SecretKey)
	session, err := s.sessionRepo.GetSessionByToken(ctx, refreshTokenHash)

	if err != nil || session == nil || session.ID != sessionID {
		return nil, "", "", ErrSessionExpired
	}

	account, err := s.accountService.GetAccountByID(ctx, session.AccountID)
	if err != nil {
		return nil, "", "", ErrUserNotFound
	}

	newAccessToken, err := s.tokenService.GenerateAccessToken(account.ID, account.Email, account.IsAdmin)
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
		session.ID,
		newRefreshTokenHash,
		newExpiresAt,
	)
	if err != nil {
		return nil, "", "", err
	}

	return account, newAccessToken, newRefreshToken, nil
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

func (s *Service) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]SessionResponse, error) {
	sessions, err := s.sessionRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var response []SessionResponse
	for _, session := range sessions {
		response = append(response, SessionResponse{
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
	return s.sessionRepo.RevokeSession(ctx, userID, sessionID)
}

func (s *Service) RevokeAllOtherSessions(
	ctx context.Context, userID uuid.UUID, currentSessionID uuid.UUID,
) (int, error) {
	sessions, err := s.sessionRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return 0, err
	}

	revokedCount := 0
	for _, session := range sessions {
		if session.ID != currentSessionID {
			err = s.sessionRepo.RevokeSession(ctx, userID, session.ID)
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
