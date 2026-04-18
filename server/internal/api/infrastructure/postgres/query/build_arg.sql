-- name: CreateBuildArg :one
INSERT INTO build_args (build_id, name, value)
VALUES (@build_id, @name, @value)
RETURNING *;

-- name: GetBuildArgs :many
SELECT ba.name, ba.value
FROM build_args ba
WHERE ba.build_id = $1;
