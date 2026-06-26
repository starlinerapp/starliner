-- name: GetDeploymentById :one
SELECT deployments.*
FROM deployments
WHERE deployments.id = @deployment_id
  AND deployments.deleted_at IS NULL;

-- name: GetDeploymentWithNamespace :one
SELECT deployments.*,
  environments.namespace
FROM deployments
  INNER JOIN environments ON deployments.environment_id = environments.id
WHERE deployments.id = $1;

-- name: GetEnvironmentDeploymentByName :one
SELECT deployments.*
FROM deployments
WHERE deployments.name = $1
  AND environment_id = $2
  AND deleted_at IS NULL;

-- name: UpdateDeploymentStatus :exec
UPDATE
  deployments
SET status = @status::deployment_status
WHERE id = @id
  AND deleted_at IS NULL;

-- name: GetDeploymentStatusLogs :one
SELECT d.status_logs,
  d.rollout_status
FROM deployments d
WHERE d.id = @deployment_id
  AND d.deleted_at IS NULL;

-- name: SetDeploymentStatusLogs :exec
UPDATE
  deployments
SET status_logs = @logs,
  rollout_status = @rollout_status
WHERE id = @deployment_id;

-- name: SoftDeleteDeployment :exec
UPDATE
  deployments
SET deleted_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: SoftDeleteDeploymentsByEnvironmentId :exec
UPDATE
  deployments
SET deleted_at = NOW()
WHERE environment_id = $1
  AND deleted_at IS NULL;

-- name: GetDeploymentsWithKubeconfig :many
SELECT deployments.*,
  c.kubeconfig,
  environments.namespace,
  c.id AS cluster_id,
  c.provisioning_id,
  c.organization_id
FROM deployments
  INNER JOIN environments ON deployments.environment_id = environments.id
  INNER JOIN projects ON environments.project_id = projects.id
  INNER JOIN clusters c ON c.id = projects.cluster_id
WHERE deployments.deleted_at IS NULL;
