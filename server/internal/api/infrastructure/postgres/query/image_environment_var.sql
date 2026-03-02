-- name: GetImageEnvironmentVars :many
SELECT ev.name, ev.value
FROM image_environment_vars ev
WHERE ev.deployment_id = $1;