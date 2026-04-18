-- name: CreateBuild :one
INSERT INTO builds (
    deployment_id,
    source
) VALUES (
    $1,
          $2
)
RETURNING *;

-- name: UpdateBuildInformation :exec
UPDATE builds
SET
    status = $1,
    commit_hash = $2,
    logs = $3
WHERE id = $4;

-- name: GetBuildLogs :one
SELECT
    b.logs
FROM builds b
INNER JOIN deployments d ON d.id = b.deployment_id
INNER JOIN environments e ON d.environment_id = e.id
INNER JOIN projects ON e.project_id = projects.id
INNER JOIN teams ON teams.id = projects.team_id
INNER JOIN team_members ON team_members.team_id = teams.id
WHERE b.id = @build_id
AND team_members.user_id = @user_id;

-- name: GetEnvironmentGitDeploymentBuilds :many
SELECT
    b.id as build_id,
    d.id as deployment_id,
    d.name as deployment_name,
    b.commit_hash,
    b.source,
    b.status,
    gd.url,
    gd.project_path,
    gd.dockerfile_path,
    b.created_at
FROM deployments d
INNER JOIN git_deployments gd ON gd.deployment_id = d.id
INNER JOIN builds b ON d.id = b.deployment_id
INNER JOIN environments e ON d.environment_id = e.id
WHERE environment_id = $1
ORDER BY b.created_at DESC;

-- name: GetLatestBuildIdByDeploymentId :one
SELECT b.id
FROM builds b
WHERE b.deployment_id = $1
ORDER BY b.created_at DESC
LIMIT 1;
