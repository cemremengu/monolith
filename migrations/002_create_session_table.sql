-- +goose Up
CREATE TABLE session
(
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    session_id   TEXT                                   NOT NULL UNIQUE,
    token_hash   TEXT                                   NOT NULL UNIQUE,
    account_id   UUID                                   NOT NULL REFERENCES account (id) ON DELETE CASCADE,
    device_info  TEXT,
    ip_address   INET,
    expires_at   TIMESTAMPTZ                            NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW()              NOT NULL,
    last_used_at TIMESTAMPTZ DEFAULT NOW()              NOT NULL,
    revoked_at   TIMESTAMPTZ
);

CREATE INDEX idx_session_account_id ON session (account_id);
CREATE INDEX idx_session_expires_at ON session (expires_at);
CREATE INDEX idx_session_token_hash ON session (token_hash);
CREATE INDEX idx_session_id ON session (session_id);
CREATE INDEX idx_session_last_used_at ON session (last_used_at);

-- +goose Down
DROP TABLE session;
