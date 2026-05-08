-- +goose Up
ALTER TABLE teams
  DROP COLUMN name,
  ADD CONSTRAINT teams_organization_id_slug_key UNIQUE (organization_id, slug);

-- +goose Down
ALTER TABLE teams
  ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '',
  DROP CONSTRAINT teams_organization_id_slug_key;

