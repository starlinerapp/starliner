-- +goose Up
ALTER TABLE clusters
    ADD COLUMN logs text;

-- +goose Down
ALTER TABLE clusters
    DROP COLUMN logs;

