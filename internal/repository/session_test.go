package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"monolith/internal/database"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SessionRepositoryTestSuite struct {
	suite.Suite
	db   *database.DB
	repo *SessionRepository
	ctx  context.Context
}

func (suite *SessionRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Skip tests if DATABASE_URL is not set (for CI/CD environments without DB)
	if testing.Short() {
		suite.T().Skip("Skipping database tests in short mode")
	}

	var err error
	suite.db, err = database.New()
	if err != nil {
		suite.T().Skipf("Failed to connect to database: %v", err)
	}

	suite.repo = NewSessionRepository(suite.db)
}

func (suite *SessionRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *SessionRepositoryTestSuite) SetupTest() {
	// Clean up tables before each test
	_, err := suite.db.Pool.Exec(suite.ctx, "DELETE FROM session")
	suite.Require().NoError(err)
	_, err = suite.db.Pool.Exec(suite.ctx, "DELETE FROM account WHERE email LIKE 'test%'")
	suite.Require().NoError(err)
}

func (suite *SessionRepositoryTestSuite) TestGenerateSessionID() {
	sessionID1, err := suite.repo.GenerateSessionID()
	suite.Require().NoError(err)
	suite.Assert().NotEmpty(sessionID1)
	suite.Assert().Len(sessionID1, sessionIDLength*2) // hex encoding doubles the length

	// Generate another ID to ensure uniqueness
	sessionID2, err := suite.repo.GenerateSessionID()
	suite.Require().NoError(err)
	suite.Assert().NotEqual(sessionID1, sessionID2)
}

func (suite *SessionRepositoryTestSuite) TestCreateSession() {
	sessionID, err := suite.repo.GenerateSessionID()
	suite.Require().NoError(err)

	tokenHash := suite.hashToken("test-token")
	accountID := suite.createTestAccount()
	deviceInfo := "Test Device"
	ipAddress := "192.168.1.1"
	expiresAt := time.Now().Add(24 * time.Hour)

	err = suite.repo.CreateSession(
		suite.ctx,
		sessionID,
		tokenHash,
		accountID,
		deviceInfo,
		ipAddress,
		expiresAt,
	)

	suite.Assert().NoError(err)

	// Verify session was created
	session, err := suite.repo.GetSessionByToken(suite.ctx, tokenHash)
	suite.Require().NoError(err)
	suite.Assert().Equal(sessionID, session.SessionID)
	suite.Assert().Equal(tokenHash, session.TokenHash)
	suite.Assert().Equal(accountID, session.AccountID)
	suite.Assert().Equal(deviceInfo, session.DeviceInfo)
	suite.Assert().Equal(session.IPAddress, deviceInfo)
	suite.Assert().WithinDuration(expiresAt, session.ExpiresAt, time.Second)
}

func (suite *SessionRepositoryTestSuite) TestGetSessionByToken_ValidToken() {
	// Create a session first
	sessionID, tokenHash, accountID := suite.createTestSession()

	session, err := suite.repo.GetSessionByToken(suite.ctx, tokenHash)

	suite.Require().NoError(err)
	suite.Assert().Equal(sessionID, session.SessionID)
	suite.Assert().Equal(tokenHash, session.TokenHash)
	suite.Assert().Equal(accountID, session.AccountID)
	suite.Assert().Nil(session.RevokedAt)
}

func (suite *SessionRepositoryTestSuite) TestGetSessionByToken_InvalidToken() {
	invalidTokenHash := suite.hashToken("invalid-token")

	session, err := suite.repo.GetSessionByToken(suite.ctx, invalidTokenHash)

	suite.Assert().Error(err)
	suite.Assert().Nil(session)
}

func (suite *SessionRepositoryTestSuite) TestGetSessionByToken_ExpiredToken() {
	sessionID, err := suite.repo.GenerateSessionID()
	suite.Require().NoError(err)

	tokenHash := suite.hashToken("expired-token")
	accountID := suite.createTestAccount()
	// Set expiration time in the past
	expiresAt := time.Now().Add(-1 * time.Hour)

	err = suite.repo.CreateSession(
		suite.ctx,
		sessionID,
		tokenHash,
		accountID,
		"Test Device",
		"192.168.1.1",
		expiresAt,
	)
	suite.Require().NoError(err)

	session, err := suite.repo.GetSessionByToken(suite.ctx, tokenHash)

	suite.Assert().Error(err)
	suite.Assert().Nil(session)
}

func (suite *SessionRepositoryTestSuite) TestGetSessionByToken_RevokedToken() {
	sessionID, tokenHash, _ := suite.createTestSession()

	// Revoke the session
	err := suite.repo.RevokeSession(suite.ctx, sessionID)
	suite.Require().NoError(err)

	session, err := suite.repo.GetSessionByToken(suite.ctx, tokenHash)

	suite.Assert().Error(err)
	suite.Assert().Nil(session)
}

func (suite *SessionRepositoryTestSuite) TestGetSessionByTokenWithTimeout_WithinTimeout() {
	_, tokenHash, _ := suite.createTestSession()

	sessionTimeout := 1 * time.Hour
	session, err := suite.repo.GetSessionByTokenWithTimeout(suite.ctx, tokenHash, sessionTimeout)

	suite.Require().NoError(err)
	suite.Assert().NotNil(session)
}

func (suite *SessionRepositoryTestSuite) TestGetSessionByTokenWithTimeout_ExceedsTimeout() {
	sessionID, err := suite.repo.GenerateSessionID()
	suite.Require().NoError(err)

	tokenHash := suite.hashToken("timeout-token")
	accountID := suite.createTestAccount()
	expiresAt := time.Now().Add(24 * time.Hour)

	// Manually insert session with old created_at time
	query := `
		INSERT INTO session (session_id, token_hash, account_id, device_info, ip_address, expires_at, created_at, rotated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
	`
	oldCreatedAt := time.Now().Add(-2 * time.Hour)
	_, err = suite.db.Pool.Exec(
		suite.ctx,
		query,
		sessionID,
		tokenHash,
		accountID,
		"Test Device",
		"192.168.1.1",
		expiresAt,
		oldCreatedAt,
	)
	suite.Require().NoError(err)

	sessionTimeout := 30 * time.Minute
	session, err := suite.repo.GetSessionByTokenWithTimeout(suite.ctx, tokenHash, sessionTimeout)

	suite.Assert().Error(err)
	suite.Assert().Nil(session)
}

func (suite *SessionRepositoryTestSuite) TestGetSessionByTokenWithTimeout_NoTimeout() {
	_, tokenHash, _ := suite.createTestSession()

	// No timeout (0 duration)
	session, err := suite.repo.GetSessionByTokenWithTimeout(suite.ctx, tokenHash, 0)

	suite.Require().NoError(err)
	suite.Assert().NotNil(session)
}

func (suite *SessionRepositoryTestSuite) TestUpdateSessionToken() {
	sessionID, oldTokenHash, _ := suite.createTestSession()

	newTokenHash := suite.hashToken("new-token")
	newExpiresAt := time.Now().Add(48 * time.Hour)

	err := suite.repo.UpdateSessionToken(suite.ctx, sessionID, newTokenHash, newExpiresAt)
	suite.Require().NoError(err)

	// Old token should not work
	session, err := suite.repo.GetSessionByToken(suite.ctx, oldTokenHash)
	suite.Assert().Error(err)
	suite.Assert().Nil(session)

	// New token should work
	session, err = suite.repo.GetSessionByToken(suite.ctx, newTokenHash)
	suite.Require().NoError(err)
	suite.Assert().Equal(sessionID, session.SessionID)
	suite.Assert().Equal(newTokenHash, session.TokenHash)
	suite.Assert().WithinDuration(newExpiresAt, session.ExpiresAt, time.Second)
}

func (suite *SessionRepositoryTestSuite) TestRevokeSession() {
	sessionID, tokenHash, _ := suite.createTestSession()

	err := suite.repo.RevokeSession(suite.ctx, sessionID)
	suite.Require().NoError(err)

	// Session should no longer be accessible
	session, err := suite.repo.GetSessionByToken(suite.ctx, tokenHash)
	suite.Assert().Error(err)
	suite.Assert().Nil(session)
}

func (suite *SessionRepositoryTestSuite) TestRevokeAllUserSessions() {
	accountID := suite.createTestAccount()

	// Create multiple sessions for the same user
	_, token1Hash := suite.createTestSessionForUser(accountID)
	_, token2Hash := suite.createTestSessionForUser(accountID)

	// Create a session for a different user
	otherAccountID := suite.createTestAccount()
	_, otherTokenHash := suite.createTestSessionForUser(otherAccountID)

	err := suite.repo.RevokeAllUserSessions(suite.ctx, accountID)
	suite.Require().NoError(err)

	// Both sessions for the target user should be revoked
	session1, err := suite.repo.GetSessionByToken(suite.ctx, token1Hash)
	suite.Assert().Error(err)
	suite.Assert().Nil(session1)

	session2, err := suite.repo.GetSessionByToken(suite.ctx, token2Hash)
	suite.Assert().Error(err)
	suite.Assert().Nil(session2)

	// Other user's session should still be valid
	otherSession, err := suite.repo.GetSessionByToken(suite.ctx, otherTokenHash)
	suite.Require().NoError(err)
	suite.Assert().NotNil(otherSession)
}

func (suite *SessionRepositoryTestSuite) TestGetUserSessions() {
	accountID := suite.createTestAccount()

	// Create multiple sessions for the user
	suite.createTestSessionForUser(accountID)
	suite.createTestSessionForUser(accountID)

	// Create a session for a different user
	otherAccountID := suite.createTestAccount()
	suite.createTestSessionForUser(otherAccountID)

	sessions, err := suite.repo.GetUserSessions(suite.ctx, accountID)

	suite.Require().NoError(err)
	suite.Assert().Len(sessions, 2)
	for _, session := range sessions {
		suite.Assert().Equal(accountID, session.AccountID)
		suite.Assert().Nil(session.RevokedAt)
	}
}

func (suite *SessionRepositoryTestSuite) TestGetUserSessions_ExcludesExpiredAndRevoked() {
	accountID := suite.createTestAccount()

	// Create a valid session
	_, _ = suite.createTestSessionForUser(accountID)

	// Create an expired session
	expiredSessionID, err := suite.repo.GenerateSessionID()
	suite.Require().NoError(err)
	expiredTokenHash := suite.hashToken("expired-token")
	err = suite.repo.CreateSession(
		suite.ctx,
		expiredSessionID,
		expiredTokenHash,
		accountID,
		"Expired Device",
		"192.168.1.2",
		time.Now().Add(-1*time.Hour), // Expired
	)
	suite.Require().NoError(err)

	// Create a revoked session
	revokedSessionID, _ := suite.createTestSessionForUser(accountID)
	err = suite.repo.RevokeSession(suite.ctx, revokedSessionID)
	suite.Require().NoError(err)

	sessions, err := suite.repo.GetUserSessions(suite.ctx, accountID)

	suite.Require().NoError(err)
	suite.Assert().Len(sessions, 1)
	suite.Assert().Equal(accountID, sessions[0].AccountID)
}

func (suite *SessionRepositoryTestSuite) TestCleanupExpiredSessions() {
	accountID := suite.createTestAccount()

	// Create a valid session
	_, validTokenHash := suite.createTestSessionForUser(accountID)

	// Create an expired session
	expiredSessionID, err := suite.repo.GenerateSessionID()
	suite.Require().NoError(err)
	expiredTokenHash := suite.hashToken("expired-token")
	err = suite.repo.CreateSession(
		suite.ctx,
		expiredSessionID,
		expiredTokenHash,
		accountID,
		"Expired Device",
		"192.168.1.2",
		time.Now().Add(-1*time.Hour), // Expired
	)
	suite.Require().NoError(err)

	// Create an old revoked session
	oldRevokedSessionID, err := suite.repo.GenerateSessionID()
	suite.Require().NoError(err)
	oldRevokedTokenHash := suite.hashToken("old-revoked-token")
	err = suite.repo.CreateSession(
		suite.ctx,
		oldRevokedSessionID,
		oldRevokedTokenHash,
		accountID,
		"Old Device",
		"192.168.1.3",
		time.Now().Add(24*time.Hour),
	)
	suite.Require().NoError(err)

	// Manually set revoked_at to be old
	_, err = suite.db.Pool.Exec(
		suite.ctx,
		"UPDATE session SET revoked_at = $1 WHERE session_id = $2",
		time.Now().Add(-31*24*time.Hour), // 31 days ago
		oldRevokedSessionID,
	)
	suite.Require().NoError(err)

	err = suite.repo.CleanupExpiredSessions(suite.ctx)
	suite.Require().NoError(err)

	// Valid session should still exist
	validSession, err := suite.repo.GetSessionByToken(suite.ctx, validTokenHash)
	suite.Require().NoError(err)
	suite.Assert().NotNil(validSession)

	// Expired and old revoked sessions should be deleted
	// Check by counting remaining sessions
	var count int
	err = suite.db.Pool.QueryRow(suite.ctx, "SELECT COUNT(*) FROM session").Scan(&count)
	suite.Require().NoError(err)
	suite.Assert().Equal(1, count) // Only the valid session should remain
}

// Helper methods

func (suite *SessionRepositoryTestSuite) createTestAccount() uuid.UUID {
	accountID := uuid.New()
	query := `
		INSERT INTO account (id, username, email, password, is_admin)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := suite.db.Pool.Exec(
		suite.ctx,
		query,
		accountID,
		"test-user-"+accountID.String()[:8],
		"test-"+accountID.String()[:8]+"@example.com",
		"$2a$12$CLuzlNmP7Bww91df6972OeKof.cFsCmKHYzfdkbExAMiAviv/PI5C",
		false,
	)
	suite.Require().NoError(err)
	return accountID
}

func (suite *SessionRepositoryTestSuite) createTestSession() (string, string, uuid.UUID) {
	sessionID, err := suite.repo.GenerateSessionID()
	suite.Require().NoError(err)

	tokenHash := suite.hashToken("test-token-" + sessionID)
	accountID := suite.createTestAccount()
	deviceInfo := "Test Device"
	ipAddress := "192.168.1.1"
	expiresAt := time.Now().Add(24 * time.Hour)

	err = suite.repo.CreateSession(
		suite.ctx,
		sessionID,
		tokenHash,
		accountID,
		deviceInfo,
		ipAddress,
		expiresAt,
	)
	suite.Require().NoError(err)

	return sessionID, tokenHash, accountID
}

func (suite *SessionRepositoryTestSuite) createTestSessionForUser(accountID uuid.UUID) (string, string) {
	sessionID, err := suite.repo.GenerateSessionID()
	suite.Require().NoError(err)

	tokenHash := suite.hashToken("test-token-" + sessionID)
	deviceInfo := "Test Device"
	ipAddress := "192.168.1.1"
	expiresAt := time.Now().Add(24 * time.Hour)

	err = suite.repo.CreateSession(
		suite.ctx,
		sessionID,
		tokenHash,
		accountID,
		deviceInfo,
		ipAddress,
		expiresAt,
	)
	suite.Require().NoError(err)

	return sessionID, tokenHash
}

func (suite *SessionRepositoryTestSuite) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func TestSessionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SessionRepositoryTestSuite))
}

// Unit tests that don't require database connection

func TestGenerateSessionID_UnitTest(t *testing.T) {
	repo := &SessionRepository{}

	sessionID, err := repo.GenerateSessionID()

	require.NoError(t, err)
	assert.NotEmpty(t, sessionID)
	assert.Len(t, sessionID, sessionIDLength*2) // hex encoding doubles the length

	// Test uniqueness
	sessionID2, err := repo.GenerateSessionID()
	require.NoError(t, err)
	assert.NotEqual(t, sessionID, sessionID2)
}

func TestSessionStruct(t *testing.T) {
	now := time.Now()
	accountID := uuid.New()
	ipAddr := "192.168.1.1"

	session := Session{
		ID:         1,
		SessionID:  "test-session-id",
		TokenHash:  "test-token-hash",
		AccountID:  accountID,
		DeviceInfo: "Test Device",
		IPAddress:  ipAddr,
		ExpiresAt:  now.Add(24 * time.Hour),
		CreatedAt:  now,
		RotatedAt:  now,
		RevokedAt:  nil,
	}

	assert.Equal(t, int64(1), session.ID)
	assert.Equal(t, "test-session-id", session.SessionID)
	assert.Equal(t, "test-token-hash", session.TokenHash)
	assert.Equal(t, accountID, session.AccountID)
	assert.Equal(t, "Test Device", session.DeviceInfo)
	assert.Equal(t, session.IPAddress, "192.168.1.1")
	assert.WithinDuration(t, now.Add(24*time.Hour), session.ExpiresAt, time.Second)
	assert.WithinDuration(t, now, session.CreatedAt, time.Second)
	assert.WithinDuration(t, now, session.RotatedAt, time.Second)
	assert.Nil(t, session.RevokedAt)
}
