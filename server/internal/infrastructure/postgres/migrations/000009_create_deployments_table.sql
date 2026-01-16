-- +goose Up
CREATE TABLE deployments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    environment_id BIGINT NOT NULL REFERENCES environments(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_deployments_updated_at
    BEFORE UPDATE ON deployments
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_deployments_updated_at ON deployments;

DROP TABLE deployments;
