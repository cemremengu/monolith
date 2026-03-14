-- name: CreateSession :one
INSERT INTO auth_session (token, prev_token, account_id, user_agent, client_ip)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSessionByToken :one
SELECT id, token, account_id, user_agent, client_ip, created_at, rotated_at, revoked_at
FROM auth_session
WHERE token = $1 OR prev_token = $2;

-- name: GetAuthContextByToken :one
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
WHERE (s.token = $1 OR s.prev_token = $2);

-- name: RotateSession :one
UPDATE auth_session
SET token = $1, prev_token = $2, rotated_at = NOW(), token_seen = FALSE, seen_at = NULL
WHERE id = $3
RETURNING *;

-- name: RevokeSession :exec
UPDATE auth_session
SET revoked_at = NOW()
WHERE id = $1 AND account_id = $2;

-- name: RevokeAllUserSessions :exec
UPDATE auth_session
SET revoked_at = NOW()
WHERE account_id = $1 AND revoked_at IS NULL;

-- name: GetSessionsByAccountID :many
SELECT id, token, account_id, user_agent, client_ip, created_at, rotated_at, revoked_at
FROM auth_session
WHERE account_id = $1 AND revoked_at IS NULL
ORDER BY rotated_at DESC;

-- name: CleanupSessions :exec
DELETE FROM auth_session
WHERE revoked_at IS NOT NULL;
