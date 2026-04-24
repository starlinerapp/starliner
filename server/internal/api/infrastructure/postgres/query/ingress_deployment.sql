-- name: CreateIngressDeployment :one
WITH new_deployment AS (
    INSERT INTO deployments (name, port, environment_id)
    VALUES (@name, @port, @environment_id)
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

-- name: UpdateIngressDeployment :one
WITH updated_ingress AS (
    UPDATE deployments
       SET port = @port
       WHERE id = @deployment_id
       RETURNING *
)
SELECT
    d.id AS deployment_id,
    d.name AS deployment_name,
    d.port AS deployment_port,
    d.status AS deployment_status,
    d.environment_id AS deployment_environment_id
FROM updated_ingress d
INNER JOIN ingress_deployments ingress_d ON d.id = ingress_d.deployment_id;

-- name: DeleteIngressPathsByDeploymentId :exec
DELETE FROM ingress_paths
WHERE deployment_id = @deployment_id;

-- name: DeleteIngressHostsByDeploymentId :exec
DELETE FROM ingress_hosts
WHERE deployment_id = @deployment_id;

-- name: GetUserEnvironmentIngressDeployments :many
SELECT
    d.id AS deployment_id,
    d.name AS deployment_name,
    d.port,
    d.status,
    d.environment_id,
    ih.id           AS host_id,
    ih.host         AS host,
    ip.id           AS path_id,
    ip.path         AS path,
    ip.path_type    AS path_type,
    svc.name        AS service_name
FROM deployments d
INNER JOIN ingress_deployments ingress_d ON d.id = ingress_d.deployment_id
INNER JOIN environments e ON d.environment_id = e.id
INNER JOIN projects ON e.project_id = projects.id
INNER JOIN teams ON projects.team_id = teams.id
INNER JOIN team_members ON team_members.team_id = teams.id
LEFT JOIN ingress_hosts ih ON ih.deployment_id = d.id
LEFT JOIN ingress_paths ip ON ip.ingress_host_id = ih.id
LEFT JOIN deployments svc ON svc.id = ip.deployment_id
WHERE d.environment_id = $1
AND team_members.user_id = $2
ORDER BY d.id DESC;

-- name: GetEnvironmentIngressDeployments :many
SELECT
    d.id AS deployment_id,
    d.name AS deployment_name,
    d.port,
    d.status,
    d.environment_id,
    ih.id           AS host_id,
    ih.host         AS host,
    ip.id           AS path_id,
    ip.path         AS path,
    ip.path_type    AS path_type,
    svc.name        AS service_name
FROM deployments d
         INNER JOIN ingress_deployments ingress_d ON d.id = ingress_d.deployment_id
         INNER JOIN environments e ON d.environment_id = e.id
         LEFT JOIN ingress_hosts ih ON ih.deployment_id = d.id
         LEFT JOIN ingress_paths ip ON ip.ingress_host_id = ih.id
         LEFT JOIN deployments svc ON svc.id = ip.deployment_id
WHERE d.environment_id = $1
ORDER BY d.id DESC;

-- name: GetEnvironmentIngressDeploymentsByName :many
SELECT
    d.id AS deployment_id,
    d.name AS deployment_name,
    d.port,
    d.status,
    d.environment_id,
    ih.id           AS host_id,
    ih.host         AS host,
    ip.id           AS path_id,
    ip.path         AS path,
    ip.path_type    AS path_type,
    svc.name        AS service_name
FROM deployments d
         INNER JOIN ingress_deployments ingress_d ON d.id = ingress_d.deployment_id
         LEFT JOIN ingress_hosts ih ON ih.deployment_id = d.id
         LEFT JOIN ingress_paths ip ON ip.ingress_host_id = ih.id
         LEFT JOIN deployments svc ON svc.id = ip.deployment_id
WHERE d.environment_id = $1
AND d.name = $2
ORDER BY d.id DESC;

-- name: GetIngressHostByName :one
SELECT i.host, i.deployment_id
FROM ingress_hosts i
WHERE i.host = $1;