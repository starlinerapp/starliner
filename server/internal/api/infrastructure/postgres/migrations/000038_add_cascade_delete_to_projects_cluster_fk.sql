-- +goose Up
ALTER TABLE projects
  DROP CONSTRAINT IF EXISTS projects_cluster_id_fkey;

ALTER TABLE projects
  ADD CONSTRAINT projects_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES clusters (id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE projects
  DROP CONSTRAINT IF EXISTS projects_cluster_id_fkey;

ALTER TABLE projects
  ADD CONSTRAINT projects_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES clusters (id);

