-- name: CreateProject :one
INSERT INTO projects (
    name,
    organization_id
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetOrganizationProjects :many
SELECT
    projects.id as id,
    projects.name as name,
    projects.organization_id as organization_id
FROM projects
INNER JOIN organizations ON projects.organization_id = organizations.id
WHERE projects.organization_id = $1;
