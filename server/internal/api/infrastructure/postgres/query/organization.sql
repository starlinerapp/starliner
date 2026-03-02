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

-- name: UpsertProvisioningCredential :exec
INSERT INTO provisioning_credentials (
    organization_id,
    provider,
    secret
) VALUES (
  $1,
  $2,
  $3
)
ON CONFLICT (organization_id, provider)
DO UPDATE SET
  secret = EXCLUDED.secret;

-- name: GetOrganizationProvisioningCredential :one
SELECT
    pc.organization_id,
    pc.secret
FROM provisioning_credentials pc
WHERE organization_id = $1
  AND provider = $2;