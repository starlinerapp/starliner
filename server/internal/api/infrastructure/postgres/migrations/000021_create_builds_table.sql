-- +goose Up
CREATE TYPE build_status AS ENUM (
  'queued', 'building', 'success', 'failure'
);

CREATE TABLE builds (
  id BIGSERIAL PRIMARY KEY, deployment_id BIGINT REFERENCES deployments (id) ON DELETE SET NULL, status build_status NOT NULL DEFAULT 'queued', logs TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_builds_updated_at
  BEFORE UPDATE ON builds
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_builds_updated_at ON builds;

DROP TABLE builds;

DROP TYPE IF EXISTS build_status;

