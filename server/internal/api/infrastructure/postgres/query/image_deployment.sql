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

-- name: CreateDeploymentVolume :one
INSERT INTO deployment_volumes (deployment_id, volume_size_mib, mount_path)
VALUES (@deployment_id, @volume_size_mib, @mount_path)
RETURNING *;

-- name: GetDeploymentVolume :one
SELECT *
FROM deployment_volumes
WHERE deployment_id = @deployment_id AND deleted_at IS NULL;

-- name: SoftDeleteDeploymentVolume :exec
UPDATE deployment_volumes
SET deleted_at = NOW()
WHERE deployment_id = @deployment_id AND deleted_at IS NULL;

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

-- name: GetUserEnvironmentImageDeployments :many
SELECT
    d.id AS deployment_id,
    d.name AS service_name,
    d.port,
    d.status,
    d.environment_id,
    img_d.name AS image_name,
    img_d.tag,
    dv.volume_size_mib,
    dv.mount_path
FROM deployments d
INNER JOIN image_deployments img_d ON d.id = img_d.deployment_id
LEFT JOIN deployment_volumes dv ON d.id = dv.deployment_id AND dv.deleted_at IS NULL
INNER JOIN environments e ON d.environment_id = e.id
INNER JOIN projects ON e.project_id = projects.id
INNER JOIN teams ON projects.team_id = teams.id
INNER JOIN team_members ON team_members.team_id = teams.id
WHERE environment_id = $1
AND team_members.user_id = $2
ORDER BY d.id DESC;

-- name: GetEnvironmentImageDeployments :many
SELECT
    d.id AS deployment_id,
    d.name AS service_name,
    d.port,
    d.status,
    d.environment_id,
    img_d.name AS image_name,
    img_d.tag,
    dv.volume_size_mib,
    dv.mount_path
FROM deployments d
         INNER JOIN image_deployments img_d ON d.id = img_d.deployment_id
         LEFT JOIN deployment_volumes dv ON d.id = dv.deployment_id AND dv.deleted_at IS NULL
         INNER JOIN environments e ON d.environment_id = e.id
WHERE environment_id = $1
ORDER BY d.id DESC;