-- name: GetPreviewEnvironment :one
SELECT environments.*, preview_environments.github_repository_id, preview_environments.pr_number
FROM environments
  INNER JOIN preview_environments ON preview_environments.environment_id = environments.id
WHERE preview_environments.github_repository_id = $1
  AND preview_environments.pr_number = $2;

