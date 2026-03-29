-- name: CreateGithubApp :one
INSERT INTO github_apps (
     installation_id, organization_id
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetOrganizationGithubApp :one
SELECT *
FROM github_apps
WHERE organization_id = $1;