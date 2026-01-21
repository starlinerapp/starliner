-- name: CreateCluster :one
INSERT INTO clusters (
    name,
    organization_id
) VALUES (
    $1,
    $2
 )
RETURNING *;

-- name: GetUserCluster :one
SELECT c.id, c.name, c.ipv4_address, c.public_key, c.private_key, c.organization_id, c.status, c.provisioning_id
FROM clusters c
LEFT JOIN organizations o ON c.organization_id = o.id
WHERE o.owner_id = $1
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
    clusters.organization_id as organization_id
FROM clusters
WHERE clusters.organization_id = $1;

-- name: GetDeploymentCluster :one
SELECT
    clusters.*
FROM clusters
INNER JOIN organizations ON organizations.id = clusters.organization_id
INNER JOIN projects ON projects.organization_id = organizations.id
INNER JOIN environments ON projects.id = environments.project_id
INNER JOIN deployments ON environments.id = deployments.environment_id
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