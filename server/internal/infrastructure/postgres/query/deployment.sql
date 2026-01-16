-- name: CreateDeployment :one
INSERT INTO deployments (
    name,
    environment_id
) VALUES (
    $1,
    $2
)
RETURNING *;
