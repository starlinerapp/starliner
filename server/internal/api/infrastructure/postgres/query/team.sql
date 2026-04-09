-- name: CreateTeam :one
INSERT INTO teams (
    name,
    slug,
    organization_id
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = $1;

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
SELECT * FROM teams
WHERE slug = $1 AND organization_id = $2;

-- name: AddTeamMember :exec
INSERT INTO team_members (
    team_id,
    user_id
) VALUES (
    $1,
    $2
)
ON CONFLICT (team_id, user_id) DO NOTHING;

-- name: RemoveTeamMember :exec
DELETE FROM team_members
WHERE team_members.team_id = $1
    AND team_members.user_id = $2;

-- name: GetTeamById :one
SELECT * FROM teams WHERE teams.id = $1;

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
        SELECT 1 FROM team_members WHERE team_members.team_id = $1
    );