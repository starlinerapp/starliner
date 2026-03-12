-- name: CreateProject :one
INSERT INTO projects (
    name,
    organization_id,
    cluster_id
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetProject :many
SELECT
    projects.id as id,
    projects.name as name,
    projects.organization_id as organization_id,
    projects.cluster_id as cluster_id,
    environments.id as environment_id,
    environments.name as environment_name,
    environments.slug as environment_slug
FROM projects
INNER JOIN organizations ON projects.organization_id = organizations.id
INNER JOIN environments ON projects.id = environments.project_id
WHERE projects.id = $1
AND organizations.owner_id = $2;

-- name: GetOrganizationProjects :many
SELECT
    projects.id as id,
    projects.name as name,
    projects.organization_id as organization_id,
    environments.id as environment_id,
    environments.name as environment_name,
    environments.slug as environment_slug,
    environments.created_at as created_at
FROM projects
INNER JOIN organizations ON projects.organization_id = organizations.id
INNER JOIN environments ON projects.id = environments.project_id
WHERE projects.organization_id = $1
ORDER BY environments.created_at;

-- name: DeleteProject :exec
DELETE
FROM projects p
USING organizations o
WHERE p.organization_id = o.id
    AND p.id = $1
    AND o.owner_id = $2;

-- name: GetProjectCluster :one
SELECT c.id, c.name
FROM projects p
INNER JOIN clusters c ON p.cluster_id = c.id
INNER JOIN organizations o ON o.id = p.organization_id
WHERE p.id = $1
    AND o.owner_id = $2;