-- name: CreateCluster :one
INSERT INTO clusters (
    name,
    ipv4_address,
    public_key,
    private_key_ref,
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