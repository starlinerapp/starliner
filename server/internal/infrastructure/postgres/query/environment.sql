-- name: CreateEnvironment :one
INSERT INTO environments (
    name,
    slug,
    project_id
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;