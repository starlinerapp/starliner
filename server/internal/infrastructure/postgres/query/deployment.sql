-- name: CreateDeployment :one
INSERT INTO deployments (
    name,
    environment_id
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetEnvironmentDeployments :many
SELECT
    deployments.*
FROM deployments
INNER JOIN environments ON deployments.environment_id = environments.id
INNER JOIN projects ON environments.project_id = projects.id
INNER JOIN organizations ON organizations.id = projects.organization_id
INNER JOIN users ON users.id = organizations.owner_id
WHERE environment_id = $1
AND users.id = $2;