-- name: CreateOrganization :one
INSERT INTO organizations (
  name, slug, owner_id)
VALUES (
  $1, $2, $3)
RETURNING *;

-- name: GetOrganization :one
SELECT *
FROM organizations
WHERE id = $1;

-- name: GetUserOrganizations :many
SELECT organizations.*
FROM organizations
  INNER JOIN organization_members ON organization_members.organization_id = organizations.id
WHERE organization_members.user_id = $1;

-- name: UpsertProvisioningCredential :exec
INSERT INTO provisioning_credentials (
  organization_id, provider, secret)
VALUES (
  $1, $2, $3)
ON CONFLICT (
  organization_id, provider)
  DO UPDATE SET
    secret = EXCLUDED.secret;

-- name: GetOrganizationProvisioningCredential :one
SELECT pc.organization_id, pc.provider, pc.secret
FROM provisioning_credentials pc
WHERE organization_id = $1
  AND provider = $2;

-- name: AddOrganizationMember :exec
INSERT INTO organization_members (
  organization_id, user_id)
VALUES (
  $1, $2)
ON CONFLICT (
  organization_id, user_id)
  DO NOTHING;

-- name: RemoveOrganizationMember :exec
DELETE FROM organization_members
WHERE organization_members.organization_id = $1
  AND organization_members.user_id = $2;

-- name: CreateOrganizationInvite :one
WITH new_invite AS (
  INSERT INTO organization_invites (
    organization_id, email, expires_at)
  VALUES (
    $1, $2, $3)
RETURNING *)
SELECT new_invite.*, organizations.name AS organization_name
FROM new_invite
  INNER JOIN organizations ON organizations.id = new_invite.organization_id;

-- name: GetOrganizationInviteById :one
SELECT organization_invites.*, organizations.name AS organization_name
FROM organization_invites
  INNER JOIN organizations ON organizations.id = organization_invites.organization_id
WHERE organization_invites.id = $1;

-- name: GetOrganizationMembers :many
SELECT users.id, users.better_auth_id
FROM users
  INNER JOIN organization_members ON organization_members.user_id = users.id
WHERE organization_members.organization_id = $1;

