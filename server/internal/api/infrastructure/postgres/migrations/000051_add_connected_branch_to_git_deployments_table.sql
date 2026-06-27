-- +goose Up
ALTER TABLE git_deployments
  ADD COLUMN connected_branch TEXT NOT NULL DEFAULT '';

UPDATE
  git_deployments gd
SET connected_branch = e.connected_branch
FROM deployments d
  INNER JOIN environments e ON d.environment_id = e.id
WHERE gd.deployment_id = d.id;

-- +goose Down
ALTER TABLE git_deployments
  DROP COLUMN connected_branch;
