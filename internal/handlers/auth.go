package handlers

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"monolith/internal/auth"
	"monolith/internal/config"
	"monolith/internal/database"
	"monolith/internal/models"
	"monolith/internal/repository"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const productionEnv = "production"

type AuthHandler struct {
	db                *database.DB
	tokenService      *auth.TokenService
	sessionRepository *repository.SessionRepository
	jwtConfig         *config.JWTConfig
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{
		db:                db,
		tokenService:      auth.NewTokenService(),
		sessionRepository: repository.NewSessionRepository(db),
		jwtConfig:         config.NewJWTConfig(),
	}
}

func (h *AuthHandler) getDeviceInfo(c echo.Context) string {
	userAgent := c.Request().Header.Get("User-Agent")
	if userAgent == "" {
		return "Unknown Device"
	}

	// Simple device detection - you could use a library like ua-parser for more sophisticated detection
	switch {
	case strings.Contains(userAgent, "iPhone"):
		return "iPhone"
	case strings.Contains(userAgent, "Android"):
		return "Android Device"
	case strings.Contains(userAgent, "Mobile"):
		return "Mobile Device"
	case strings.Contains(userAgent, "Chrome"):
		return "Chrome Browser"
	case strings.Contains(userAgent, "Firefox"):
		return "Firefox Browser"
	case strings.Contains(userAgent, "Safari"):
		return "Safari Browser"
	default:
		return "Desktop Browser"
	}
}

func (h *AuthHandler) getClientIP(c echo.Context) string {
	// Check for X-Forwarded-For header first (common in load balancers/proxies)
	if xff := c.Request().Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check for X-Real-IP header (common in nginx)
	if xri := c.Request().Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	return c.RealIP()
}

func (h *AuthHandler) generateAndSetTokens(c echo.Context, userID, email string, isAdmin bool) error {
	accessToken, err := h.tokenService.GenerateAccessToken(userID, email, isAdmin)
	if err != nil {
		return err
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(userID)
	if err != nil {
		return err
	}

	// Create session
	sessionID, err := h.sessionRepository.GenerateSessionID()
	if err != nil {
		return err
	}

	refreshTokenHash := auth.HashToken(refreshToken)
	accountID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	deviceInfo := h.getDeviceInfo(c)
	ipAddress := h.getClientIP(c)
	expiresAt := time.Now().Add(h.tokenService.RefreshTokenDuration())

	err = h.sessionRepository.CreateSession(
		c.Request().Context(),
		sessionID,
		refreshTokenHash,
		accountID,
		deviceInfo,
		ipAddress,
		expiresAt,
	)
	if err != nil {
		return err
	}

	// Set cookies
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.tokenService.AccessTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(accessCookie)

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.tokenService.RefreshTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(refreshCookie)

	// Store session ID in a cookie for session management
	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.tokenService.RefreshTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(sessionCookie)

	return nil
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Password == "" || len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be at least 8 characters"})
	}

	var existingUser models.User
	err := pgxscan.Get(context.Background(), h.db.Pool, &existingUser, `
		SELECT id FROM account WHERE email = $1 OR username = $2
	`, req.Email, req.Username)
	if err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "User already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	var user models.User
	err = pgxscan.Get(context.Background(), h.db.Pool, &user, `
		INSERT INTO account (username, email, name, password, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, NOW(), NOW()) 
		RETURNING id, username, email, name, is_admin, language, theme, timezone, 
		          last_seen_at, is_disabled, created_at, updated_at
	`, req.Username, req.Email, req.Name, string(hashedPassword))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	if tokenErr := h.generateAndSetTokens(c, user.ID, user.Email, user.IsAdmin); tokenErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate tokens"})
	}

	response := struct {
		User models.User `json:"user"`
	}{
		User: user,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	var user models.User
	err := pgxscan.Get(context.Background(), h.db.Pool, &user, `
		SELECT id, username, email, name, password, is_admin, language, theme, timezone, 
		       last_seen_at, is_disabled, created_at, updated_at
		FROM account 
		WHERE (email = $1 OR username = $1) AND is_disabled = FALSE
	`, req.Login)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	_, err = h.db.Pool.Exec(context.Background(), `
		UPDATE account SET last_seen_at = NOW() WHERE id = $1
	`, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update last seen"})
	}

	if tokenErr := h.generateAndSetTokens(c, user.ID, user.Email, user.IsAdmin); tokenErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate tokens"})
	}

	response := struct {
		User models.User `json:"user"`
	}{
		User: user,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	var user models.User
	err := pgxscan.Get(context.Background(), h.db.Pool, &user, `
		SELECT id, username, email, name, is_admin, language, theme, timezone, 
		       last_seen_at, is_disabled, created_at, updated_at
		FROM account 
		WHERE id = $1 AND is_disabled = FALSE
	`, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// Revoke session if it exists
	if sessionCookie, cookieErr := c.Cookie("session_id"); cookieErr == nil {
		_ = h.sessionRepository.RevokeSession(c.Request().Context(), sessionCookie.Value)
	}

	// Clear all auth cookies
	cookies := []string{"access_token", "refresh_token", "session_id"}
	for _, cookieName := range cookies {
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    "",
			HttpOnly: true,
			Secure:   os.Getenv("ENV") == productionEnv,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1, // Delete cookie
			Path:     "/",
		}
		c.SetCookie(cookie)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	refreshTokenCookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Refresh token not found"})
	}

	sessionCookie, err := c.Cookie("session_id")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Session not found"})
	}

	// Validate refresh token
	claims, err := h.tokenService.ValidateRefreshToken(refreshTokenCookie.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid refresh token"})
	}

	// Check if session exists, refresh token matches, and session hasn't timed out
	refreshTokenHash := auth.HashToken(refreshTokenCookie.Value)
	session, err := h.sessionRepository.GetSessionByTokenWithTimeout(
		c.Request().Context(),
		refreshTokenHash,
		h.jwtConfig.SessionTimeout,
	)
	if err != nil || session == nil || session.SessionID != sessionCookie.Value {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Session expired or invalid"})
	}

	// Get user information
	var user models.User
	err = pgxscan.Get(context.Background(), h.db.Pool, &user, `
		SELECT id, username, email, name, is_admin, language, theme, timezone, 
		       last_seen_at, is_disabled, created_at, updated_at
		FROM account 
		WHERE id = $1 AND is_disabled = FALSE
	`, claims.UserID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
	}

	// Generate new tokens
	newAccessToken, err := h.tokenService.GenerateAccessToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate access token"})
	}

	newRefreshToken, err := h.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate refresh token"})
	}

	// Update session with new refresh token (instead of creating a new session)
	newRefreshTokenHash := auth.HashToken(newRefreshToken)
	newExpiresAt := time.Now().Add(h.tokenService.RefreshTokenDuration())
	err = h.sessionRepository.UpdateSessionToken(
		c.Request().Context(),
		session.SessionID,
		newRefreshTokenHash,
		newExpiresAt,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update session"})
	}

	// Set new cookies
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.tokenService.AccessTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(accessCookie)

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == productionEnv,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.tokenService.RefreshTokenDuration().Seconds()),
		Path:     "/",
	}
	c.SetCookie(refreshCookie)

	return c.JSON(http.StatusOK, map[string]string{"message": "Tokens refreshed successfully"})
}

func (h *AuthHandler) GetSessions(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}
	accountID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	sessions, err := h.sessionRepository.GetUserSessions(c.Request().Context(), accountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch sessions"})
	}

	// Get current session ID to mark it
	var currentSessionID string
	if sessionCookie, cookieErr := c.Cookie("session_id"); cookieErr == nil {
		currentSessionID = sessionCookie.Value
	}

	type SessionResponse struct {
		SessionID  string    `json:"sessionId"`
		DeviceInfo string    `json:"deviceInfo"`
		IPAddress  string    `json:"ipAddress"`
		CreatedAt  time.Time `json:"createdAt"`
		RotatedAt  time.Time `json:"rotatedAt"`
		IsCurrent  bool      `json:"isCurrent"`
	}

	var response []SessionResponse
	for _, session := range sessions {
		response = append(response, SessionResponse{
			SessionID:  session.SessionID,
			DeviceInfo: session.DeviceInfo,
			IPAddress:  session.IPAddress,
			CreatedAt:  session.CreatedAt,
			RotatedAt:  session.RotatedAt,
			IsCurrent:  session.SessionID == currentSessionID,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RevokeSession(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}
	sessionID := c.Param("sessionId")

	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Session ID is required"})
	}

	// Verify the session belongs to the user
	accountID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	sessions, err := h.sessionRepository.GetUserSessions(c.Request().Context(), accountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify session"})
	}

	sessionFound := false
	for _, session := range sessions {
		if session.SessionID == sessionID {
			sessionFound = true
			break
		}
	}

	if !sessionFound {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Session not found"})
	}

	err = h.sessionRepository.RevokeSession(c.Request().Context(), sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to revoke session"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Session revoked successfully"})
}

func (h *AuthHandler) RevokeAllOtherSessions(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}
	accountID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	// Get current session ID to preserve it
	currentSessionID := ""
	if sessionCookie, cookieErr := c.Cookie("session_id"); cookieErr == nil {
		currentSessionID = sessionCookie.Value
	}

	// Get all user sessions
	sessions, err := h.sessionRepository.GetUserSessions(c.Request().Context(), accountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch sessions"})
	}

	// Revoke all sessions except the current one
	revokedCount := 0
	for _, session := range sessions {
		if session.SessionID != currentSessionID {
			err = h.sessionRepository.RevokeSession(c.Request().Context(), session.SessionID)
			if err == nil {
				revokedCount++
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message":      "Other sessions revoked successfully",
		"revokedCount": revokedCount,
	})
}
