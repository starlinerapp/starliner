-- +goose Up
ALTER TABLE database_deployments
    ADD COLUMN database TEXT;

-- +goose Down
ALTER TABLE database_deployments
    DROP COLUMN database;