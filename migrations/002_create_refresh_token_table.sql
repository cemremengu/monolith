-- +goose Up
CREATE TABLE refresh_token
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    token_hash TEXT                                   NOT NULL UNIQUE,
    account_id UUID                                   NOT NULL REFERENCES account (id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ                            NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()              NOT NULL,
    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_refresh_token_account_id ON refresh_token (account_id);
CREATE INDEX idx_refresh_token_expires_at ON refresh_token (expires_at);
CREATE INDEX idx_refresh_token_hash ON refresh_token (token_hash);

-- +goose Down
DROP TABLE refresh_token;