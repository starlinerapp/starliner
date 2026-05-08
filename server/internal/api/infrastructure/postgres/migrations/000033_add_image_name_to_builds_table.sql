-- +goose Up
ALTER TABLE builds
  ADD COLUMN image_name TEXT;

-- +goose Down
ALTER TABLE builds
  DROP COLUMN image_name;

