-- +goose Up
ALTER TABLE database_deployments
    ADD COLUMN DATABASE TEXT;

-- +goose Down
ALTER TABLE database_deployments
    DROP COLUMN DATABASE;

