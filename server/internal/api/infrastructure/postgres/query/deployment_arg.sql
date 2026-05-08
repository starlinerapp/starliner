-- name: CreateGitDeploymentArg :one
INSERT INTO git_deployment_args (
  deployment_id, name, value)
VALUES (
  @deployment_id, @name, @value)
RETURNING *;

-- name: GetGitDeploymentArgs :many
SELECT da.name, da.value
FROM git_deployment_args da
WHERE da.deployment_id = $1;

-- name: DeleteArgsByDeploymentId :exec
DELETE FROM git_deployment_args
WHERE deployment_id = $1;

