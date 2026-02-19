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

-- +goose Down
DROP TRIGGER IF EXISTS trigger_ingress_deployments_updated_at ON ingress_deployments;

DROP TABLE ingress_deployments;
