-- +goose Up
ALTER TABLE projects
    ADD COLUMN preview_environments_enabled BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE projects
    DROP COLUMN preview_environments_enabled;
