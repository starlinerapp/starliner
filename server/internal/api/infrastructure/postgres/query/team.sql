-- name: CreateTeam :one
INSERT INTO teams (
  slug, organization_id)
VALUES (
  $1, $2)
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1;

-- name: GetUserTeams :many
SELECT teams.*
FROM teams
  INNER JOIN team_members ON team_members.team_id = teams.id
WHERE teams.organization_id = $1
  AND team_members.user_id = $2;

-- name: GetTeamMembers :many
SELECT users.id, users.better_auth_id
FROM team_members
  INNER JOIN users ON team_members.user_id = users.id
WHERE team_members.team_id = $1;

-- name: GetTeamBySlug :one
SELECT *
FROM teams
WHERE slug = $1
  AND organization_id = $2;

-- name: AddTeamMember :exec
INSERT INTO team_members (
  team_id, user_id)
VALUES (
  $1, $2)
ON CONFLICT (
  team_id, user_id)
  DO NOTHING;

-- name: RemoveTeamMember :exec
DELETE FROM team_members
WHERE team_members.team_id = $1
  AND team_members.user_id = $2;

-- name: GetTeamById :one
SELECT *
FROM teams
WHERE teams.id = $1;

-- name: FindTeamByIdAndUserId :one
SELECT teams.*
FROM teams
  INNER JOIN organization_members ON organization_members.organization_id = teams.organization_id
WHERE teams.id = $1
  AND organization_members.user_id = $2;

-- name: DeleteTeamIfEmpty :exec
DELETE FROM teams
WHERE id = $1
  AND NOT EXISTS (
    SELECT 1
    FROM team_members
    WHERE team_members.team_id = $1);

-- name: AssignRepoToTeam :exec
INSERT INTO team_repositories (
  team_id, github_repo_id, repo_name, github_app_id)
VALUES (
  $1, $2, $3, $4)
ON CONFLICT (
  team_id, github_repo_id)
  DO NOTHING;

-- name: UnassignRepoFromTeam :exec
DELETE FROM team_repositories
WHERE team_id = $1
  AND github_repo_id = $2;

-- name: GetTeamRepositories :many
SELECT github_repo_id, repo_name
FROM team_repositories
WHERE team_id = $1;

-- name: GetTeamsByRepoAndOrg :many
SELECT teams.*
FROM teams
  INNER JOIN team_repositories ON team_repositories.team_id = teams.id
WHERE teams.organization_id = $1
  AND team_repositories.github_repo_id = $2;

-- name: GetTeamClusters :many
SELECT clusters.*
FROM clusters
  INNER JOIN team_clusters ON clusters.id = team_clusters.cluster_id
WHERE team_clusters.team_id = $1;

-- name: AssignTeamCluster :exec
INSERT INTO team_clusters (
  team_id, cluster_id)
VALUES (
  $1, $2)
ON CONFLICT (
  team_id, cluster_id)
  DO NOTHING;

-- name: UnassignTeamCluster :exec
DELETE FROM team_clusters
WHERE team_clusters.team_id = $1
  AND team_clusters.cluster_id = $2;

-- name: GetTeamCluster :one
SELECT clusters.*
FROM team_clusters
  INNER JOIN clusters ON clusters.id = team_clusters.cluster_id
WHERE team_clusters.team_id = $1
  AND team_clusters.cluster_id = $2;

