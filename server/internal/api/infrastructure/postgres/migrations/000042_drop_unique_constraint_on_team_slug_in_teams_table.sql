-- +goose Up
ALTER TABLE teams
    DROP CONSTRAINT teams_slug_key;

-- +goose Down
ALTER TABLE teams
    ADD CONSTRAINT teams_slug_key UNIQUE (slug);