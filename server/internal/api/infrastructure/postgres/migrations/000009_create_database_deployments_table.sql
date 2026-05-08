-- +goose Up
CREATE TABLE deployments (
  id BIGSERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, port VARCHAR(255) NOT NULL, status VARCHAR(255), environment_id BIGINT NOT NULL REFERENCES environments (id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_deployments_updated_at
  BEFORE UPDATE ON deployments
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column ();

CREATE TABLE database_deployments (
  deployment_id BIGINT PRIMARY KEY REFERENCES deployments (id) ON DELETE CASCADE, username VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_database_deployments_updated_at
  BEFORE UPDATE ON database_deployments
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_database_deployments_updated_at ON database_deployments;

DROP TRIGGER IF EXISTS trigger_deployments_updated_at ON deployments;

DROP TABLE database_deployments;

DROP TABLE deployments;

