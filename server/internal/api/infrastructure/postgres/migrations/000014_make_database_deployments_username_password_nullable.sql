-- +goose Up
ALTER TABLE database_deployments
    ALTER COLUMN username DROP NOT NULL,
    ALTER COLUMN PASSWORD DROP NOT NULL;

-- +goose Down
UPDATE
    database_deployments
SET
    username = COALESCE(username, ''),
    password = COALESCE(password, '')
WHERE
    username IS NULL
    OR PASSWORD IS NULL;

ALTER TABLE database_deployments
    ALTER COLUMN username SET NOT NULL,
    ALTER COLUMN PASSWORD SET NOT NULL;

