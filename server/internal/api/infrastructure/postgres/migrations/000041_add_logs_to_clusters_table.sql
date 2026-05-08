-- +goose Up
ALTER TABLE clusters
  ADD COLUMN logs TEXT;

-- +goose Down
ALTER TABLE clusters
  DROP COLUMN logs;

