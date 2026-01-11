-- name: CreateDeployment :exec
INSERT INTO deployments (
    name,
    environment_id
) VALUES (
    $1,
    $2
);
