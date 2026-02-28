-- name: CreateIngressDeployment :one
WITH new_deployment AS (
    INSERT INTO deployments (name, port, status, environment_id)
    VALUES (@name, @port, @status, @environment_id)
    RETURNING *
),
new_ingress_deployment AS (
    INSERT INTO ingress_deployments (deployment_id)
    SELECT id FROM new_deployment
    RETURNING *
)
SELECT
    d.id AS deployment_id,
    d.name AS deployment_name,
    d.port AS deployment_port,
    d.status AS deployment_status,
    d.environment_id AS deployment_environment_id
FROM new_deployment d
INNER JOIN new_ingress_deployment ingress_d ON d.id = ingress_d.deployment_id;

-- name: GetEnvironmentIngressDeployments :many
SELECT
    d.id AS deployment_id,
    d.name AS service_name,
    d.port,
    d.status,
    d.environment_id
FROM deployments d
INNER JOIN ingress_deployments ingress_d ON d.id = ingress_d.deployment_id
INNER JOIN environments e ON d.environment_id = e.id
INNER JOIN projects ON e.project_id = projects.id
INNER JOIN organizations ON organizations.id = projects.organization_id
INNER JOIN users ON users.id = organizations.owner_id
WHERE environment_id = $1
AND users.id = $2;