-- name: UserExists :one
SELECT id FROM account WHERE email = $1 OR username = $2;

-- name: GetAccountByLogin :one
SELECT id, username, email, name, password, is_admin, language, theme, timezone,
       last_seen_at, status, created_at, updated_at
FROM account
WHERE (email = $1 OR username = $1) AND status = 'active';

-- name: GetAccountByID :one
SELECT id, username, email, name, is_admin, language, theme, timezone,
       last_seen_at, status, created_at, updated_at
FROM account
WHERE id = $1 AND status = 'active';

-- name: GetAccountByIDWithPassword :one
SELECT id, password
FROM account
WHERE id = $1 AND status = 'active';

-- name: RegisterAccount :one
INSERT INTO account (username, email, name, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, username, email, name, is_admin, language, theme, timezone,
          last_seen_at, status, created_at, updated_at;

-- name: CreateAccount :one
INSERT INTO account (username, name, email, password, is_admin, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
RETURNING id, username, email, name, avatar, is_admin, language, theme, timezone,
          last_seen_at, status, created_at, updated_at;

-- name: CreateInvitedAccount :one
INSERT INTO account (username, email, name, is_admin, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, 'pending', NOW(), NOW())
RETURNING id, username, email, name, avatar, is_admin, language, theme, timezone,
          last_seen_at, status, created_at, updated_at;

-- name: UpdateLastSeen :exec
UPDATE account SET last_seen_at = NOW() WHERE id = $1;

-- name: UpdatePreferences :one
UPDATE account
SET language = $1, theme = $2, timezone = $3, updated_at = NOW()
WHERE id = $4 AND status = 'active'
RETURNING id, username, email, name, is_admin, language, theme, timezone,
          last_seen_at, status, created_at, updated_at;

-- name: UpdateAccount :one
UPDATE account
SET username = $1, name = $2, email = $3, updated_at = NOW()
WHERE id = $4 AND status = 'active'
RETURNING id, username, email, name, avatar, is_admin, language, theme, timezone,
          last_seen_at, status, created_at, updated_at;

-- name: UpdatePassword :exec
UPDATE account SET password = $1, updated_at = NOW() WHERE id = $2;

-- name: DisableAccount :exec
UPDATE account SET status = 'disabled', updated_at = NOW() WHERE id = $1;

-- name: EnableAccount :exec
UPDATE account SET status = 'active', updated_at = NOW() WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM account WHERE id = $1;

-- name: GetAccounts :many
SELECT id, username, email, name, avatar, is_admin, language, theme, timezone,
       last_seen_at, status, created_at, updated_at
FROM account
ORDER BY created_at DESC;

-- name: GetAccount :one
SELECT id, username, email, name, avatar, is_admin, language, theme, timezone,
       last_seen_at, status, created_at, updated_at
FROM account
WHERE id = $1 AND status = 'active';
