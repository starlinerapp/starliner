-- +goose Up
ALTER TABLE deployments
  ADD COLUMN rollout_status TEXT NOT NULL DEFAULT 'pending';

ALTER TABLE deployments
  ADD CONSTRAINT deployments_rollout_status_check CHECK (rollout_status IN ('pending', 'success', 'failure'));

UPDATE
  deployments
SET rollout_status = CASE WHEN status_logs_complete
  AND COALESCE(status_logs, '')
  LIKE '%has failed.%' THEN
    'failure'
  WHEN status_logs_complete
  AND COALESCE(status_logs, '')
  LIKE '%is complete.%' THEN
    'success'
  ELSE
    'pending'
  END;

-- +goose Down
ALTER TABLE deployments
  DROP CONSTRAINT IF EXISTS deployments_rollout_status_check;

ALTER TABLE deployments
  DROP COLUMN IF EXISTS rollout_status;
