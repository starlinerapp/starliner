-- +goose Up
ALTER TABLE deployments
  ADD COLUMN status_logs TEXT,
  ADD COLUMN status_logs_complete BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE deployments
  DROP COLUMN IF EXISTS status_logs_complete,
  DROP COLUMN IF EXISTS status_logs;

