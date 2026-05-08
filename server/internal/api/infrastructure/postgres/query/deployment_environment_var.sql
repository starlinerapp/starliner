-- name: GetDeploymentEnvironmentVars :many
SELECT ev.name, ev.value
FROM deployment_environment_vars ev
WHERE ev.deployment_id = $1;

-- name: DeleteEnvVarsByDeploymentId :exec
DELETE FROM deployment_environment_vars
WHERE deployment_id = $1;

