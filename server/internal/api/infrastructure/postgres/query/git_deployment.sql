-- name: CreateGitDeployment :one
WITH new_deployment AS (
    INSERT INTO deployments (name, port, environment_id)
    VALUES ($1, $2, $3)
    RETURNING *
),
new_git_deployment AS (
    INSERT INTO git_deployments (deployment_id, url, project_path, dockerfile_path)
    SELECT id, $4, $5, $6 FROM new_deployment
    RETURNING *
)
SELECT
    d.id AS deployment_id,
    d.name,
    d.port,
    d.environment_id,
    gd.url,
    gd.dockerfile_path,
    gd.project_path
FROM new_deployment d
INNER JOIN new_git_deployment gd ON d.id = gd.deployment_id;

-- name: GetEnvironmentGitDeployments :many
SELECT
    d.id AS deployment_id,
    d.name,
    d.port,
    d.status,
    d.environment_id,
    gd.url,
    gd.project_path,
    gd.dockerfile_path
FROM deployments d
INNER JOIN git_deployments gd ON d.id = gd.deployment_id
INNER JOIN environments ON d.environment_id = environments.id
INNER JOIN projects ON environments.project_id = projects.id
INNER JOIN organizations ON organizations.id = projects.organization_id
INNER JOIN users ON users.id = organizations.owner_id
WHERE environment_id = @environment_id
  AND users.id = @user_id
ORDER BY d.id DESC;