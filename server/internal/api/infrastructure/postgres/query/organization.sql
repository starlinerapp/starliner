-- name: CreateOrganization :one
INSERT INTO organizations (
    name,
    slug,
    owner_id
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetOrganization :one
SELECT *
FROM organizations
WHERE id = $1;

-- name: GetUserOrganizations :many
SELECT *
FROM organizations
WHERE owner_id = $1;