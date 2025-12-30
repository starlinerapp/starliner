-- name: CreateCluster :one
INSERT INTO clusters (
    name,
    ipv4_address,
    public_key,
    private_key,
    organization_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
 )
RETURNING *;

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
    clusters.ipv4_address as ipv4_address,
    clusters.public_key as public_key,
    clusters.private_key as private_key,
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