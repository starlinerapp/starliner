-- +goose Up
ALTER TABLE projects RENAME COLUMN organization_id TO team_id;

ALTER TABLE projects
  DROP CONSTRAINT projects_organization_id_fkey;

ALTER TABLE projects
  ADD CONSTRAINT projects_team_id_fkey FOREIGN KEY (team_id) REFERENCES teams (id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE projects
  DROP CONSTRAINT projects_team_id_fkey;

ALTER TABLE projects
  ADD CONSTRAINT projects_organization_id_fkey FOREIGN KEY (team_id) REFERENCES organizations (id) ON DELETE CASCADE;

ALTER TABLE projects RENAME COLUMN team_id TO organization_id;

