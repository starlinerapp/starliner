-- name: CreateGitDeployment :one
WITH new_deployment AS (
INSERT INTO deployments (name, port, environment_id)
        VALUES ($1, $2, $3)
    RETURNING
        *
), new_git_deployment AS (
INSERT INTO git_deployments (deployment_id, url, project_path, dockerfile_path)
    SELECT
        id,
        $4,
        $5,
        $6
    FROM
        new_deployment
    RETURNING
        *
)
SELECT
    d.id AS deployment_id,
    d.name,
    d.port,
    d.environment_id,
    gd.url,
    gd.dockerfile_path,
    gd.project_path
FROM
    new_deployment d
    INNER JOIN new_git_deployment gd ON d.id = gd.deployment_id;

-- name: UpdateGitDeployment :one
WITH updated_deployment AS (
    UPDATE
        deployments
    SET
        port = @port
    WHERE
        id = @deployment_id
    RETURNING
        *
),
updated_git_deployment AS (
    UPDATE
        git_deployments
    SET
        project_path = @project_path,
        dockerfile_path = @dockerfile_path
    WHERE
        deployment_id = @deployment_id
    RETURNING
        *
)
SELECT
    d.id AS deployment_id,
    d.status,
    d.name AS service_name,
    git_d.url,
    git_d.project_path,
    git_d.dockerfile_path,
    d.port,
    d.environment_id
FROM
    updated_deployment d
    INNER JOIN updated_git_deployment git_d ON d.id = git_d.deployment_id;

-- name: GetUserEnvironmentGitDeployments :many
SELECT
    d.id AS deployment_id,
    d.name,
    d.port,
    d.status,
    d.environment_id,
    gd.url,
    gd.project_path,
    gd.dockerfile_path
FROM
    deployments d
    INNER JOIN git_deployments gd ON d.id = gd.deployment_id
    INNER JOIN environments ON d.environment_id = environments.id
    INNER JOIN projects ON environments.project_id = projects.id
    INNER JOIN teams ON projects.team_id = teams.id
    INNER JOIN team_members ON team_members.team_id = teams.id
WHERE
    environment_id = @environment_id
    AND team_members.user_id = @user_id
ORDER BY
    d.id DESC;

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
FROM
    deployments d
    INNER JOIN git_deployments gd ON d.id = gd.deployment_id
    INNER JOIN environments ON d.environment_id = environments.id
WHERE
    environment_id = @environment_id
ORDER BY
    d.id DESC;

-- name: GetGitDeploymentsByRepositoryUrl :many
SELECT
    d.id AS deployment_id,
    d.name,
    d.port,
    d.status,
    d.environment_id,
    gd.url,
    gd.project_path,
    gd.dockerfile_path
FROM
    deployments d
    INNER JOIN git_deployments gd ON d.id = gd.deployment_id
WHERE
    gd.url = @repository_url;

