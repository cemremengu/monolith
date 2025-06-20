-- +goose Up
CREATE TABLE account
(
    id           UUID                     DEFAULT gen_random_uuid() PRIMARY KEY,
    username     TEXT                                   NOT NULL,
    email        TEXT UNIQUE                            NOT NULL,
    name         TEXT,
    password     TEXT,
    avatar       TEXT,
    salt         TEXT,
    rands        TEXT,
    is_admin     BOOLEAN                  DEFAULT FALSE,
    language     TEXT,
    theme        TEXT,
    timezone     TEXT,
    last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_disabled  boolean                  DEFAULT FALSE not null,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose StatementBegin
INSERT INTO public.account (id, username, email, name, password, salt, rands, is_admin, language, theme, timezone, last_seen_at, is_disabled, created_at, updated_at) VALUES ('61856bb8-f32c-44da-bf51-216047be674c', 'admin', 'admin@ttgint.com', null, '2ddd92fabc31d83403793e700321ff5310680ee4cce275b76d516d650a49072caf3f1237955400e40f29febf818290fd0e36', 'BeFyGJvd3l', null, TRUE, null, null, null, NOW(), FALSE, NOW(), NOW());
-- +goose StatementEnd

-- +goose Down
DROP TABLE account;
