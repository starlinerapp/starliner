-- name: CreateBuild :one
INSERT INTO builds (
    deployment_id
) VALUES (
    $1
)
RETURNING *;

-- name: UpdateBuildInformation :exec
UPDATE builds
SET
    status = $1,
    logs = $2
WHERE id = $3;