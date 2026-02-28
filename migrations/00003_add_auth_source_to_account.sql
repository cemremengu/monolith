-- +goose Up
ALTER TABLE account ADD COLUMN auth_source TEXT NOT NULL DEFAULT 'local';

-- +goose Down
ALTER TABLE account DROP COLUMN auth_source;
