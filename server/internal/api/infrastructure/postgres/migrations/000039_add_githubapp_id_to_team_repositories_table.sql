-- +goose Up
DELETE FROM team_repositories;

ALTER TABLE team_repositories
    ADD COLUMN github_app_id bigint NOT NULL;

ALTER TABLE team_repositories
    ADD CONSTRAINT team_repositories_github_app_id_fkey FOREIGN KEY (github_app_id) REFERENCES github_apps (id) ON DELETE CASCADE;

CREATE INDEX idx_team_repositories_github_app_id ON team_repositories (github_app_id);

-- +goose Down
DROP INDEX IF EXISTS idx_team_repositories_github_app_id;

ALTER TABLE team_repositories
    DROP CONSTRAINT IF EXISTS team_repositories_github_app_id_fkey;

ALTER TABLE team_repositories
    DROP COLUMN IF EXISTS github_app_id;

