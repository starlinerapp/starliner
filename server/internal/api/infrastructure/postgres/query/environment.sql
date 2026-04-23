-- name: CreateEnvironment :one
INSERT INTO environments (
    name,
    slug,
    namespace,
    project_id
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: CreateEnvironmentWithConnectedBranch :one
INSERT INTO environments (
    name,
    slug,
    namespace,
    project_id,
    connected_branch
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: CreatePreviewEnvironment :one
WITH new_environment AS (
    INSERT INTO environments (
        name,
        slug,
        namespace,
        project_id,
        connected_branch
    ) VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) RETURNING *
),
new_preview_environments AS (
    INSERT INTO preview_environments (
        environment_id,
        github_repository_id,
        pr_number
    )
   SELECT id, $6, $7 FROM new_environment
   RETURNING *
)
SELECT
    e.id,
    e.name,
    e.slug,
    e.project_id,
    e.namespace,
    e.connected_branch,
    pe.github_repository_id,
    pe.pr_number
FROM new_environment e
INNER JOIN new_preview_environments pe ON e.id = pe.environment_id;

-- name: GetEnvironmentById :one
SELECT *
FROM environments
WHERE environments.id = $1;

-- name: GetEnvironmentProject :one
SELECT p.*
FROM projects p
INNER JOIN environments e on p.id = e.project_id
WHERE e.id = $1;

-- name: GetEnvironmentCluster :one
SELECT clusters.*
FROM environments
INNER JOIN projects ON projects.id = environments.project_id
INNER JOIN clusters ON projects.cluster_id = clusters.id
WHERE environments.id = $1;

-- name: GetEnvironmentAuthorizedUsers :many
SELECT tm.user_id
FROM environments
INNER JOIN projects p ON p.id = environments.project_id
INNER JOIN teams t ON t.id = p.team_id
INNER JOIN team_members tm ON tm.team_id = t.id
WHERE environments.id = $1;

-- name: GetEnvironmentBranch :one
SELECT e.connected_branch
FROM environments e
WHERE e.id = $1;

-- name: UpdateEnvironmentBranch :exec
UPDATE environments
SET connected_branch = $1
WHERE id = $2;