-- +goose Up
CREATE TABLE team_repositories (
  team_id BIGINT NOT NULL REFERENCES teams (id) ON DELETE CASCADE, github_repo_id BIGINT NOT NULL, repo_name VARCHAR(255) NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), PRIMARY KEY (team_id, github_repo_id)
);

CREATE TRIGGER trigger_team_repositories_updated_at
  BEFORE UPDATE ON team_repositories
  FOR EACH ROW
  EXECUTE PROCEDURE update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_team_repositories_updated_at ON team_repositories;

DROP TABLE team_repositories;

