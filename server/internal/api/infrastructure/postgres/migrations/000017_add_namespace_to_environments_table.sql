-- +goose Up
ALTER TABLE environments
  ADD COLUMN namespace TEXT NOT NULL DEFAULT 'default';

-- +goose Down
ALTER TABLE environments
  DROP COLUMN namespace;

