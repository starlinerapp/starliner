-- +goose Up
CREATE TABLE teams (
  id BIGSERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, slug VARCHAR(255) NOT NULL UNIQUE, organization_id BIGINT NOT NULL REFERENCES organizations (id) ON DELETE CASCADE, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_builds_teams
  BEFORE UPDATE ON teams
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column ();

CREATE TABLE team_members (
  team_id BIGINT NOT NULL REFERENCES teams (id) ON DELETE CASCADE, user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), PRIMARY KEY (team_id, user_id)
);

CREATE TRIGGER trigger_team_members_updated_at
  BEFORE UPDATE ON team_members
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_team_members_updated_at ON team_members;

DROP TABLE team_members;

DROP TRIGGER IF EXISTS trigger_builds_teams ON teams;

DROP TABLE teams;

