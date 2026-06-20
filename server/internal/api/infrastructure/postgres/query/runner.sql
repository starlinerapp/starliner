-- name: CreateRunner :one
INSERT INTO runners (
  organization_id)
VALUES (
  $1)
RETURNING *;

-- name: CreateRunnerRegistrationToken :one
INSERT INTO runner_registration_tokens (
  organization_id,
  token_hash,
  expires_at,
  runner_id)
VALUES (
  $1,
  $2,
  $3,
  $4)
RETURNING *;

-- name: UseRunnerRegistrationToken :one
UPDATE
  runner_registration_tokens
SET used_at = NOW()
WHERE token_hash = $1
  AND used_at IS NULL
  AND expires_at > NOW()
RETURNING *;

-- name: UpdateRunnerOnRegistration :one
UPDATE
  runners
SET name = $2,
  labels = $3,
  max_concurrent_jobs = $4
WHERE id = $1
RETURNING *;

-- name: GetOrganizationRunners :many
SELECT *
FROM runners
WHERE organization_id = $1
ORDER BY created_at DESC;

-- name: GetRunnerIdByRegistrationToken :one
SELECT runner_id
FROM runner_registration_tokens
WHERE token_hash = $1
  AND runner_id IS NOT NULL;

-- name: UpdateRunnerStatus :exec
UPDATE
  runners
SET status = $2
WHERE id = $1;

-- name: GetRunnerByOrganization :one
SELECT *
FROM runners
WHERE id = $1
  AND organization_id = $2;

-- name: DeleteRunner :exec
DELETE FROM runners
WHERE id = $1
  AND organization_id = $2;
