-- name: CreateEnvironment :one
INSERT INTO environments (
    name,
    slug,
    namespace,
    project_id
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: CreateEnvironmentWithConnectedBranch :one
INSERT INTO environments (
    name,
    slug,
    namespace,
    project_id,
    connected_branch
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetEnvironmentById :one
SELECT *
FROM environments
WHERE environments.id = $1;

-- name: GetEnvironmentProject :one
SELECT p.*
FROM projects p
INNER JOIN environments e on p.id = e.project_id
WHERE e.id = $1;

-- name: GetEnvironmentCluster :one
SELECT clusters.*
FROM environments
INNER JOIN projects ON projects.id = environments.project_id
INNER JOIN clusters ON projects.cluster_id = clusters.id
WHERE environments.id = $1;

-- name: GetEnvironmentAuthorizedUsers :many
SELECT tm.user_id
FROM environments
INNER JOIN projects p ON p.id = environments.project_id
INNER JOIN teams t ON t.id = p.team_id
INNER JOIN team_members tm ON tm.team_id = t.id
WHERE environments.id = $1;

-- name: GetEnvironmentBranch :one
SELECT e.connected_branch
FROM environments e
WHERE e.id = $1;

-- name: UpdateEnvironmentBranch :exec
UPDATE environments
SET connected_branch = $1
WHERE id = $2;