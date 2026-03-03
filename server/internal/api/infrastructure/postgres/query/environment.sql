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

-- name: GetEnvironmentById :one
SELECT *
FROM environments
WHERE environments.id = $1;

-- name: GetEnvironmentCluster :one
SELECT clusters.*
FROM environments
INNER JOIN projects ON projects.id = environments.project_id
INNER JOIN clusters ON projects.cluster_id = clusters.id
WHERE environments.id = $1;

-- name: GetEnvironmentAuthorizedUsers :many
SELECT u.id
FROM environments
INNER JOIN projects p ON p.id = environments.project_id
INNER JOIN organizations o ON o.id = p.organization_id
INNER JOIN users u ON o.owner_id = u.id
WHERE environments.id = $1;
