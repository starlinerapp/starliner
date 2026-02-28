-- +goose Up
CREATE TABLE ingress_deployments (
    deployment_id BIGINT PRIMARY KEY REFERENCES deployments(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_ingress_deployments_updated_at
    BEFORE UPDATE ON ingress_deployments
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE ingress_hosts (
    id BIGSERIAL PRIMARY KEY,
    deployment_id BIGINT NOT NULL REFERENCES ingress_deployments(deployment_id) ON DELETE CASCADE,
    host TEXT NOT NULL,

    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_ingress_host UNIQUE (deployment_id, host)
);

CREATE INDEX idx_ingress_hosts_deployment_id ON ingress_hosts (deployment_id);

CREATE TRIGGER trigger_ingress_hosts_updated_at
    BEFORE UPDATE ON ingress_hosts
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE ingress_paths (
   id BIGSERIAL PRIMARY KEY,
   ingress_host_id BIGINT NOT NULL REFERENCES ingress_hosts(id) ON DELETE CASCADE,

   deployment_id BIGINT NOT NULL REFERENCES deployments(id) ON DELETE CASCADE,

   path TEXT NOT NULL,
   path_type TEXT NOT NULL,

   created_at timestamptz NOT NULL DEFAULT NOW(),
   updated_at timestamptz NOT NULL DEFAULT NOW(),

   CONSTRAINT unique_ingress_paths UNIQUE (ingress_host_id, path, path_type, deployment_id)
);

CREATE INDEX idx_ingress_paths_ingress_host_id
    ON ingress_paths (ingress_host_id);

CREATE TRIGGER trigger_ingress_paths_updated_at
    BEFORE UPDATE ON ingress_paths
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_ingress_paths_updated_at ON ingress_paths;
DROP TABLE ingress_paths;

DROP TRIGGER IF EXISTS trigger_ingress_hosts_updated_at ON ingress_hosts;
DROP TABLE ingress_hosts;

DROP TRIGGER IF EXISTS trigger_ingress_deployments_updated_at ON ingress_deployments;
DROP TABLE ingress_deployments;
