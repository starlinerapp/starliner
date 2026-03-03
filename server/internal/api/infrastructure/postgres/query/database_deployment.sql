-- name: CreateDatabaseDeployment :one
WITH new_deployment AS (
    INSERT INTO deployments (name, port, environment_id)
    VALUES ($1, $2, $3)
    RETURNING *
),
new_database_deployment AS (
    INSERT INTO database_deployments (deployment_id)
    SELECT id FROM new_deployment
    RETURNING *
)
SELECT
    d.id AS deployment_id,
    d.name,
    d.port,
    d.environment_id,
    db.username,
    db.password
FROM new_deployment d
INNER JOIN new_database_deployment db ON d.id = db.deployment_id;

-- name: UpdateDatabaseDeploymentCredentials :exec
UPDATE database_deployments
SET database = @database,
    username = @username,
    password = @password
WHERE deployment_id = @deployment_id;

-- name: GetEnvironmentDatabaseDeployments :many
SELECT
    d.id AS deployment_id,
    d.name,
    d.port,
    d.status,
    d.environment_id,
    db.database,
    db.username,
    db.password
FROM deployments d
INNER JOIN database_deployments db ON d.id = db.deployment_id
INNER JOIN environments ON d.environment_id = environments.id
INNER JOIN projects ON environments.project_id = projects.id
INNER JOIN organizations ON organizations.id = projects.organization_id
INNER JOIN users ON users.id = organizations.owner_id
WHERE environment_id = $1
AND users.id = $2
ORDER BY d.id DESC;