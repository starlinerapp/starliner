-- +goose Up
ALTER TABLE clusters
  ADD COLUMN server_type TEXT NOT NULL DEFAULT 'cx23';

-- +goose Down
ALTER TABLE clusters
  DROP COLUMN server_type;

