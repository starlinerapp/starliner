-- +goose Up
ALTER TABLE clusters
  ADD COLUMN "user" TEXT NOT NULL DEFAULT 'root';

ALTER TABLE clusters
  DROP CONSTRAINT IF EXISTS clusters_name_key;

ALTER TABLE clusters
  ADD CONSTRAINT clusters_name_organization_id_key UNIQUE (NAME, organization_id);

-- +goose Down
ALTER TABLE clusters
  DROP CONSTRAINT IF EXISTS clusters_name_organization_id_key;

ALTER TABLE clusters
  ADD CONSTRAINT clusters_name_key UNIQUE (NAME);

ALTER TABLE clusters
  DROP COLUMN IF EXISTS "user";

