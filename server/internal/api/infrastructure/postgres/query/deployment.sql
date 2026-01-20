-- name: GetUserDeployment :one
SELECT deployments.*
FROM  deployments
INNER JOIN environments ON deployments.environment_id = environments.id
INNER JOIN projects ON environments.project_id = projects.id
INNER JOIN organizations o on projects.organization_id = o.id
INNER JOIN users ON o.owner_id = users.id
WHERE deployments.id = @deployment_id
    AND users.id = @user_id;