-- +goose Up
ALTER TABLE builds
    ADD COLUMN commit_hash text,
    ADD COLUMN source text NOT NULL DEFAULT 'manual';

-- +goose Down
ALTER TABLE builds
    DROP COLUMN IF EXISTS commit_hash,
    DROP COLUMN IF EXISTS source;

