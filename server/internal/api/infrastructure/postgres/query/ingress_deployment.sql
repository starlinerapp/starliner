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

-- name: CreateIngressHost :one
INSERT INTO ingress_hosts (deployment_id, host)
VALUES (@deployment_id, @host)
ON CONFLICT (deployment_id, host) DO UPDATE SET host = EXCLUDED.host
RETURNING id, deployment_id, host;

-- name: CreateIngressPath :one
INSERT INTO ingress_paths (ingress_host_id, deployment_id, path, path_type)
VALUES (@ingress_host_id, @deployment_id, @path, @path_type)
ON CONFLICT (ingress_host_id, path, path_type, deployment_id) DO UPDATE SET path = EXCLUDED.path
RETURNING id, ingress_host_id, deployment_id, path, path_type;

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
AND users.id = $2
ORDER BY d.id DESC;