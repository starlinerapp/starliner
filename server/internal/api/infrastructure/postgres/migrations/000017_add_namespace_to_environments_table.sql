-- +goose Up
ALTER TABLE environments
    ADD COLUMN namespace text NOT NULL DEFAULT 'default';

-- +goose Down
ALTER TABLE environments
    DROP COLUMN namespace;

