-- +goose Up
CREATE TYPE deployment_status AS ENUM (
  'healthy', 'unhealthy'
);

UPDATE
  deployments
SET status = 'unhealthy'
WHERE status IS NULL
  OR status NOT IN ('healthy', 'unhealthy');

ALTER TABLE deployments
  ALTER COLUMN status TYPE deployment_status
  USING status::deployment_status;

ALTER TABLE deployments
  ALTER COLUMN status SET DEFAULT 'unhealthy',
  ALTER COLUMN status SET NOT NULL;

-- +goose Down
ALTER TABLE deployments
  ALTER COLUMN status DROP DEFAULT,
  ALTER COLUMN status DROP NOT NULL;

ALTER TABLE deployments
  ALTER COLUMN status TYPE VARCHAR(255)
  USING status::TEXT;

DROP TYPE deployment_status;

