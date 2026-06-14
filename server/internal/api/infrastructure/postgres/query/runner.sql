-- name: CreateRunner :one
INSERT INTO runners (
  organization_id)
VALUES (
  $1)
RETURNING *;

-- name: CreateRunnerRegistrationToken :one
INSERT INTO runner_registration_tokens (
  organization_id, token_hash, expires_at, runner_id)
VALUES (
  $1, $2, $3, $4)
RETURNING *;
