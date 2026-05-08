-- +goose Up
CREATE TYPE provider AS ENUM (
  'hetzner'
);

CREATE TABLE provisioning_credentials (
  id BIGSERIAL PRIMARY KEY, organization_id BIGINT NOT NULL REFERENCES organizations (id) ON DELETE CASCADE, provider provider NOT NULL, secret TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), CONSTRAINT unique_provider_per_organization UNIQUE (organization_id, provider)
);

CREATE INDEX provisioning_credentials_organization_id ON provisioning_credentials (organization_id);

CREATE TRIGGER update_provisioning_credentials_updated_at
  BEFORE UPDATE ON provisioning_credentials
  FOR EACH ROW
  EXECUTE PROCEDURE update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS update_provisioning_credentials_updated_at ON provisioning_credentials;

DROP TABLE provisioning_credentials;

DROP TYPE IF EXISTS provider;

