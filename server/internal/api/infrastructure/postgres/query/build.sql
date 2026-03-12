-- name: CreateBuild :one
INSERT INTO builds (
    deployment_id
) VALUES (
    $1
)
RETURNING *;

-- name: UpdateBuildInformation :exec
UPDATE builds
SET
    status = $1,
    logs = $2
WHERE id = $3;

-- name: GetEnvironmentGitDeploymentBuilds :many
SELECT
    b.id as build_id,
    d.id as deployment_id,
    d.name as deployment_name,
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