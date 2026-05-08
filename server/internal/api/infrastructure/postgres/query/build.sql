-- name: CreateBuild :one
INSERT INTO builds (deployment_id, source)
    VALUES ($1, $2)
RETURNING
    *;

-- name: UpdateBuildInformation :exec
UPDATE
    builds
SET
    status = $1,
    commit_hash = $2,
    image_name = $3,
    logs = $4
WHERE
    id = $5;

-- name: GetBuildLogs :one
SELECT
    b.logs
FROM
    builds b
    INNER JOIN deployments d ON d.id = b.deployment_id
    INNER JOIN environments e ON d.environment_id = e.id
    INNER JOIN projects ON e.project_id = projects.id
    INNER JOIN teams ON teams.id = projects.team_id
    INNER JOIN team_members ON team_members.team_id = teams.id
WHERE
    b.id = @build_id
    AND team_members.user_id = @user_id;

-- name: GetEnvironmentGitDeploymentBuilds :many
SELECT
    b.id AS build_id,
    d.id AS deployment_id,
    d.name AS deployment_name,
    b.image_name AS image_name,
    b.commit_hash,
    b.source,
    b.status,
    gd.url,
    gd.project_path,
    gd.dockerfile_path,
    b.created_at
FROM
    deployments d
    INNER JOIN git_deployments gd ON gd.deployment_id = d.id
    INNER JOIN builds b ON d.id = b.deployment_id
    INNER JOIN environments e ON d.environment_id = e.id
WHERE
    environment_id = $1
ORDER BY
    b.created_at DESC;

-- name: GetLatestGitDeploymentBuild :one
SELECT DISTINCT ON (d.id)
    b.id AS build_id,
    d.id AS deployment_id,
    d.name AS deployment_name,
    b.image_name AS image_name,
    b.commit_hash,
    b.source,
    b.status,
    gd.url,
    gd.project_path,
    gd.dockerfile_path,
    b.created_at
FROM
    deployments d
    INNER JOIN git_deployments gd ON gd.deployment_id = d.id
    INNER JOIN builds b ON d.id = b.deployment_id
    INNER JOIN environments e ON d.environment_id = e.id
WHERE
    d.environment_id = $1
    AND d.name = $2
ORDER BY
    d.id,
    b.created_at DESC;

