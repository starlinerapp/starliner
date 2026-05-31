-- +goose Up
ALTER TABLE deployments
  DROP COLUMN IF EXISTS status_logs_complete;

-- +goose Down
ALTER TABLE deployments
  ADD COLUMN status_logs_complete BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE
  deployments
SET status_logs_complete = CASE WHEN COALESCE(status_logs, '') != '' THEN
    TRUE
  ELSE
    FALSE
  END;

