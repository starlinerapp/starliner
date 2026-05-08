-- +goose Up
ALTER TABLE projects
  ADD COLUMN cluster_id BIGINT REFERENCES clusters (id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE projects
  DROP COLUMN cluster_id;

