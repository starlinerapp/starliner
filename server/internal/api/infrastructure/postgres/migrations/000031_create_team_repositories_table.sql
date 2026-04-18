-- +goose Up
CREATE TABLE team_repositories (
    team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    github_repo_id BIGINT NOT NULL,
    repo_name VARCHAR(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (team_id, github_repo_id)
);

-- +goose Down
DROP TABLE team_repositories;

