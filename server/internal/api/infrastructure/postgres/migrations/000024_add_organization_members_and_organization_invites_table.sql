-- +goose Up
CREATE TABLE organization_members (
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (organization_id, user_id)
);

CREATE TABLE organization_invites (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
      expires_at timestamptz NOT NULL,
      created_at timestamptz NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE organization_invites;

DROP TABLE organization_members;
