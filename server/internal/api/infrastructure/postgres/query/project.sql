-- name: CreateProject :one
INSERT INTO projects (
    name,
    team_id,
    cluster_id
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetProject :many
SELECT
    projects.id as id,
    projects.name as name,
    projects.team_id as team_id,
    projects.cluster_id as cluster_id,
    environments.id as environment_id,
    environments.name as environment_name,
    environments.slug as environment_slug
FROM projects
INNER JOIN teams ON projects.team_id = teams.id
INNER JOIN team_members ON teams.id = team_members.team_id
INNER JOIN environments ON projects.id = environments.project_id
WHERE projects.id = $1
    AND team_members.user_id = $2;

-- name: GetUserProjects :many
SELECT
    projects.id as id,
    projects.name as name,
    projects.team_id as team_id,
    environments.id as environment_id,
    environments.name as environment_name,
    environments.slug as environment_slug,
    environments.created_at as created_at
FROM projects
INNER JOIN teams ON projects.team_id = teams.id
INNER JOIN team_members ON teams.id = team_members.team_id
INNER JOIN environments ON projects.id = environments.project_id
WHERE teams.organization_id = $1
  AND team_members.user_id = $2
ORDER BY environments.created_at;

-- name: DeleteProject :exec
DELETE
FROM projects p
USING organizations o, teams t
WHERE t.organization_id = o.id
  AND p.team_id = t.id
  AND p.id = $1
  AND o.owner_id = $2;

-- name: GetProjectCluster :one
SELECT c.id, c.name
FROM projects p
INNER JOIN clusters c ON p.cluster_id = c.id
INNER JOIN teams t ON p.team_id = t.id
INNER JOIN team_members tm ON tm.team_id = t.id
WHERE p.id = $1
  AND tm.user_id = $2;

-- name: GetProjectEnvironments :many
SELECT e.*
FROM environments e
         INNER JOIN projects p ON p.id = e.project_id
         INNER JOIN teams t ON t.id = p.team_id
         INNER JOIN team_members tm ON tm.team_id = t.id
WHERE e.project_id = $1 AND tm.user_id = $2;

-- name: GetProjectProductionEnvironmentsByRepositoryUrl :many
SELECT e.*
FROM environments e
WHERE e.name = 'Production'
  AND EXISTS (
    SELECT 1
    FROM deployments d
             JOIN git_deployments g ON g.deployment_id = d.id
    WHERE d.environment_id = e.id
      AND g.url = $1
);

-- name: GetProjectPreviewEnvironmentEnabled :one
SELECT p.preview_environments_enabled
FROM projects p
    INNER JOIN teams t ON t.id = p.team_id
    INNER JOIN team_members tm ON tm.team_id = t.id
WHERE p.id = $1 AND tm.user_id = $2;

-- name: ToggleProjectPreviewEnvironmentEnabled :one
UPDATE projects p
SET preview_environments_enabled = NOT p.preview_environments_enabled
FROM teams t
         INNER JOIN team_members tm ON tm.team_id = t.id
WHERE p.team_id = t.id
  AND p.id = $1
  AND tm.user_id = $2
RETURNING p.preview_environments_enabled;
