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
    environments.slug as environment_slug
FROM projects
INNER JOIN organizations ON projects.organization_id = organizations.id
INNER JOIN environments ON projects.id = environments.project_id
WHERE projects.organization_id = $1;