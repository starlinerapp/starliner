-- +goose NO TRANSACTION
-- +goose Up
ALTER TYPE cluster_status ADD VALUE IF NOT EXISTS 'failed';

ALTER TABLE clusters
  ADD COLUMN deleted_at TIMESTAMPTZ;

CREATE INDEX idx_clusters_deleted_at ON clusters (deleted_at)
WHERE deleted_at IS NULL;

ALTER TABLE clusters
  DROP CONSTRAINT IF EXISTS clusters_name_organization_id_key;

-- +goose Down
ALTER TABLE clusters
  ADD CONSTRAINT clusters_name_organization_id_key UNIQUE (NAME, organization_id);

DROP INDEX IF EXISTS idx_clusters_deleted_at;

ALTER TABLE clusters
  DROP COLUMN IF EXISTS deleted_at;