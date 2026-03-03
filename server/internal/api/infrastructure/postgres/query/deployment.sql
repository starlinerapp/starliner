-- name: GetUserDeployment :one
SELECT deployments.*
FROM  deployments
INNER JOIN environments ON deployments.environment_id = environments.id
INNER JOIN projects ON environments.project_id = projects.id
INNER JOIN organizations o on projects.organization_id = o.id
INNER JOIN users ON o.owner_id = users.id
WHERE deployments.id = @deployment_id
    AND users.id = @user_id;

-- name: GetEnvironmentDeploymentByName :one
SELECT deployments.*
FROM deployments
INNER JOIN environments ON deployments.environment_id = environments.id
WHERE deployments.name = $1
    AND environment_id = $2;
;

-- name: UpdateDeploymentStatus :exec
UPDATE deployments
SET status = @status::deployment_status
WHERE id = @id;

-- name: DeleteDeployment :exec
DELETE
FROM deployments
where id = $1;


-- name: GetDeploymentsWithKubeconfig :many
SELECT deployments.*, c.kubeconfig
FROM deployments
INNER JOIN environments ON deployments.environment_id = environments.id
INNER JOIN projects ON environments.project_id = projects.id
INNER JOIN clusters c on c.id = projects.cluster_id;