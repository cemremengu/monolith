-- +goose Up
CREATE TABLE account
(
    id           UUID                     DEFAULT gen_random_uuid() PRIMARY KEY,
    username     TEXT                                   NOT NULL,
    email        TEXT UNIQUE                            NOT NULL,
    name         TEXT,
    password     TEXT,
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

-- +goose Down
DROP TABLE account;
