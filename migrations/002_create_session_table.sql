-- +goose Up
CREATE TABLE session
(
    id         UUID        DEFAULT gen_random_uuid() PRIMARY KEY,
    token_hash TEXT                      NOT NULL UNIQUE,
    account_id UUID                      NOT NULL REFERENCES account (id) ON DELETE CASCADE,
    user_agent TEXT,
    client_ip  TEXT,
    expires_at TIMESTAMPTZ               NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    rotated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_session_account_id ON session (account_id);
CREATE INDEX idx_session_expires_at ON session (expires_at);
CREATE INDEX idx_session_revoked_at ON session (revoked_at);

-- +goose Down
DROP TABLE session;
