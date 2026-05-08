-- name: CreateGithubApp :one
INSERT INTO github_apps (
  installation_id, organization_id)
VALUES (
  $1, $2)
RETURNING *;

-- name: GetOrganizationGithubApp :one
SELECT *
FROM github_apps
WHERE organization_id = $1;

-- name: GetEnvironmentGithubApp :one
SELECT ga.*
FROM github_apps ga
  INNER JOIN organizations o ON o.id = ga.organization_id
  INNER JOIN teams t ON t.organization_id = o.id
  INNER JOIN projects p ON p.team_id = t.id
  INNER JOIN environments e ON e.project_id = p.id
WHERE e.id = $1;

-- name: DeleteGithubAppByInstallationId :exec
DELETE FROM github_apps
WHERE installation_id = $1;

