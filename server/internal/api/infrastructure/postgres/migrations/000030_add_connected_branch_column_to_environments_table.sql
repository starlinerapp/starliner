-- +goose Up
ALTER TABLE environments
  ADD COLUMN connected_branch TEXT NOT NULL DEFAULT 'main';

-- +goose Down
ALTER TABLE environments
  DROP COLUMN connected_branch;

