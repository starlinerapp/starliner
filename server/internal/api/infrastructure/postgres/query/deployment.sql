-- name: GetUserDeployment :one
SELECT
    deployments.*
FROM
    deployments
    INNER JOIN environments ON deployments.environment_id = environments.id
    INNER JOIN projects ON environments.project_id = projects.id
    INNER JOIN teams ON projects.team_id = teams.id
    INNER JOIN team_members ON team_members.team_id = teams.id
WHERE
    deployments.id = @deployment_id
    AND team_members.user_id = @user_id;

-- name: GetDeploymentWithNamespace :one
SELECT
    deployments.*,
    environments.namespace
FROM
    deployments
    INNER JOIN environments ON deployments.environment_id = environments.id
WHERE
    deployments.id = $1;

-- name: GetEnvironmentDeploymentByName :one
SELECT
    deployments.*
FROM
    deployments
WHERE
    deployments.name = $1
    AND environment_id = $2;

-- name: UpdateDeploymentStatus :exec
UPDATE
    deployments
SET
    status = @status::deployment_status
WHERE
    id = @id;

-- name: DeleteDeployment :exec
DELETE FROM deployments
WHERE id = $1;

-- name: GetDeploymentsWithKubeconfig :many
SELECT
    deployments.*,
    c.kubeconfig,
    environments.namespace
FROM
    deployments
    INNER JOIN environments ON deployments.environment_id = environments.id
    INNER JOIN projects ON environments.project_id = projects.id
    INNER JOIN clusters c ON c.id = projects.cluster_id;

