-- +goose Up
ALTER TABLE builds
    ADD COLUMN image_name text;

-- +goose Down
ALTER TABLE builds
    DROP COLUMN image_name;

