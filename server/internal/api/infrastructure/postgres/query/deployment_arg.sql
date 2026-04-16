-- name: CreateDeploymentArg :one
INSERT INTO deployment_args (deployment_id, name, value)
VALUES (@deployment_id, @name, @value)
RETURNING *;

-- name: GetDeploymentArgs :many
SELECT da.name, da.value
FROM deployment_args da
WHERE da.deployment_id = $1;

-- name: DeleteArgsByDeploymentId :exec
DELETE FROM deployment_args
WHERE deployment_id = $1;

