-- +goose Up
ALTER TABLE organization_invites
  ADD COLUMN team_id BIGINT REFERENCES teams (id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE organization_invites
  DROP COLUMN team_id;

