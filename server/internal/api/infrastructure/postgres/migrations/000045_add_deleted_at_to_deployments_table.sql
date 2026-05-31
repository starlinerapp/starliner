-- +goose Up
ALTER TABLE deployments
  ADD COLUMN deleted_at TIMESTAMPTZ;

CREATE INDEX idx_deployments_deleted_at ON deployments (deleted_at)
WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_deployments_deleted_at;

ALTER TABLE deployments
  DROP COLUMN IF EXISTS deleted_at;

