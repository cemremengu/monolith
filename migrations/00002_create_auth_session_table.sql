-- +goose Up
CREATE TABLE auth_session
(
    id         UUID        DEFAULT gen_random_uuid() PRIMARY KEY,
    token      TEXT                      NOT NULL UNIQUE,
    prev_token TEXT                      NOT NULL UNIQUE,
    account_id UUID                      NOT NULL,
    user_agent TEXT                      NOT NULL,
    client_ip  TEXT                      NOT NULL,
    token_seen BOOLEAN     DEFAULT TRUE  NOT NULL,
    seen_at    TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    rotated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    revoked_at TIMESTAMPTZ
);

CREATE INDEX ix_auth_session_account_id ON auth_session (account_id);
CREATE INDEX ix_auth_session_revoked_at ON auth_session (revoked_at);

-- +goose Down
DROP TABLE auth_session;
