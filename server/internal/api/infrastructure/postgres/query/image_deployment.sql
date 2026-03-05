-- name: CreateImageDeployment :one
WITH new_deployment AS (
    INSERT INTO deployments (name, port, environment_id)
    VALUES (@service_name, @port, @environment_id)
    RETURNING *
),
new_image_deployment AS (
    INSERT INTO image_deployments (deployment_id, name, tag)
    SELECT id, @image_name, @tag FROM new_deployment
    RETURNING *
)
SELECT
    d.id AS deployment_id,
    d.status,
    d.name AS service_name,
    img_d.name AS image_name,
    img_d.tag AS image_tag,
    d.port,
    d.environment_id
FROM new_deployment d
INNER JOIN new_image_deployment img_d ON d.id = img_d.deployment_id;

-- name: UpdateImageDeployment :one
WITH updated_deployment AS (
    UPDATE deployments
        SET port = @port
        WHERE id = @deployment_id
        RETURNING *
),
updated_image_deployment AS (
 UPDATE image_deployments
     SET
         name = @image_name,
         tag  = @tag
     WHERE deployment_id = (SELECT id FROM updated_deployment)
     RETURNING *
)
SELECT
    d.id AS deployment_id,
    d.status,
    d.name AS service_name,
    img_d.name AS image_name,
    img_d.tag AS image_tag,
    d.port,
    d.environment_id
FROM updated_deployment d
INNER JOIN updated_image_deployment img_d ON d.id = img_d.deployment_id;

-- name: CreateDeploymentEnvVar :one
INSERT INTO deployment_environment_vars (deployment_id, name, value)
VALUES (@deployment_id, @name, @value)
RETURNING *;

-- name: GetEnvironmentImageDeployments :many
SELECT
    d.id AS deployment_id,
    d.name AS service_name,
    d.port,
    d.status,
    d.environment_id,
    img_d.name AS image_name,
    img_d.tag
FROM deployments d
INNER JOIN image_deployments img_d ON d.id = img_d.deployment_id
INNER JOIN environments e ON d.environment_id = e.id
INNER JOIN projects ON e.project_id = projects.id
INNER JOIN organizations ON organizations.id = projects.organization_id
INNER JOIN users ON users.id = organizations.owner_id
WHERE environment_id = $1
AND users.id = $2
ORDER BY d.id DESC;