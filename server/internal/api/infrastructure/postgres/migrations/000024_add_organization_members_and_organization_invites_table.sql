-- +goose Up
CREATE TABLE organization_members (
  organization_id BIGINT NOT NULL REFERENCES organizations (id) ON DELETE CASCADE, user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), PRIMARY KEY (organization_id, user_id)
);

CREATE TRIGGER trigger_organization_members_updated_at
  BEFORE UPDATE ON organization_members
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column ();

CREATE TABLE organization_invites (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (), organization_id BIGINT NOT NULL REFERENCES organizations (id) ON DELETE CASCADE, expires_at TIMESTAMPTZ NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE organization_invites;

DROP TRIGGER IF EXISTS trigger_organization_members_updated_at ON organization_members;

DROP TABLE organization_members;

