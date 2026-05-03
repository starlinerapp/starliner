-- name: CreateCluster :one
INSERT INTO clusters (
    name,
    server_type,
    organization_id
) VALUES (
    $1,
    $2,
    $3
 )
RETURNING *;

-- name: GetUserCluster :one
SELECT c.id, c.name, c.user, c.ipv4_address, c.public_key, c.private_key, c.organization_id, c.status, c.provisioning_id, c.server_type
FROM clusters c
LEFT JOIN organizations o ON c.organization_id = o.id
LEFT JOIN organization_members om ON o.id = om.organization_id
WHERE om.user_id = $1
AND c.id = $2;

-- name: GetCluster :one
SELECT *
FROM clusters
WHERE id = $1;

-- name: DeleteCluster :exec
DELETE
FROM clusters
WHERE id = $1;

-- name: GetOrganizationClusters :many
SELECT
    clusters.id as id,
    clusters.name as name,
    clusters.organization_id as organization_id,
    clusters.server_type as server_type,
    clusters.created_at as created_at
FROM clusters
WHERE clusters.organization_id = $1;

-- name: GetDeploymentCluster :one
SELECT
    clusters.*
FROM clusters
INNER JOIN projects ON projects.cluster_id = clusters.id
INNER JOIN environments ON environments.project_id = projects.id
INNER JOIN deployments ON deployments.environment_id = environments.id
WHERE deployments.id = @deployment_id;

-- name: UpdateClusterPublicPrivateKeys :exec
UPDATE clusters
SET
    public_key = $1,
    private_key = $2
WHERE id = $3;

-- name: UpdateClusterIPv4Address :exec
UPDATE clusters
SET
    ipv4_address = $1
WHERE id = $2;

-- name: UpdateClusterProvisioningId :exec
UPDATE clusters
SET
    provisioning_id = $1
WHERE id = $2;

-- name: UpdateClusterStatus :exec
UPDATE clusters
SET
    status = $1
WHERE id = $2;

-- name: UpdateClusterKubeconfig :exec
UPDATE clusters
SET
    kubeconfig = $1
where id = $2;

-- name: UpdateClusterLogs :exec
UPDATE clusters
SET
    logs = $1
WHERE id = $2;

-- name: GetUserClusterProvisioningLogs :one
SELECT
    c.logs
FROM clusters c
LEFT JOIN organizations o ON c.organization_id = o.id
LEFT JOIN organization_members om ON o.id = om.organization_id
WHERE c.id = @cluster_id
AND om.user_id = @user_id;