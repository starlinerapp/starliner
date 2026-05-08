-- name: CreateUser :one
INSERT INTO "users" (better_auth_id)
    VALUES ($1)
RETURNING
    *;

-- name: GetUserByBetterAuthId :one
SELECT
    *
FROM
    "users"
WHERE
    better_auth_id = $1;

