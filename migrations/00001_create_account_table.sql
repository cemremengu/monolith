-- +goose Up
CREATE TABLE account
(
    id           UUID        DEFAULT gen_random_uuid() PRIMARY KEY,
    username     TEXT                      NOT NULL UNIQUE,
    email        TEXT                      NOT NULL UNIQUE,
    name         TEXT,
    avatar       TEXT,
    password     TEXT,
    is_admin     BOOLEAN     DEFAULT FALSE,
    language     TEXT,
    theme        TEXT,
    timezone     TEXT,
    last_seen_at TIMESTAMPTZ DEFAULT NOW(),
    is_disabled  BOOLEAN     DEFAULT FALSE not null,
    status       TEXT        DEFAULT 'disabled'::TEXT NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX account_username_email_idx ON account (username, email);

INSERT INTO public.account (id, username, email, name, avatar, password, is_admin, language, theme, timezone, last_seen_at, status, created_at, updated_at) VALUES ('4645fd03-84ac-44ac-b26b-9178fd67de17', 'admin', 'admin@localhost.com', 'System Admin', null, '$2a$12$PRlEsSCdZr4guV.hA6efmuHPf/QKqJlvPnfKp043OfpwxDac1kSeW', true, null, null, null, now(), 'active', now(), now());

-- +goose Down
DROP TABLE account;
