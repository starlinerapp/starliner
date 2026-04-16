-- +goose Up
CREATE TABLE deployment_args (
    id BIGSERIAL PRIMARY KEY,
    deployment_id BIGINT NOT NULL REFERENCES deployments(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    value TEXT NOT NULL,

    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX deployment_args_deployment_id_name ON deployment_args(deployment_id, name);

CREATE TRIGGER trigger_deployment_args_updated_at
    BEFORE UPDATE ON deployment_args
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_deployment_args_updated_at ON deployment_args;
DROP TABLE deployment_args;

