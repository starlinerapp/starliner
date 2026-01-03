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
SELECT c.id, c.name, c.ipv4_address, c.public_key, c.private_key, c.organization_id, c.status, c.pulumi_stack_id
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

-- name: UpdateClusterPulumiStackId :exec
UPDATE clusters
SET
    pulumi_stack_id = $1
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