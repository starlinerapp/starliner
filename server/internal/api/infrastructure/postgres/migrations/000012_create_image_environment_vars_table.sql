-- +goose Up
CREATE TABLE image_environment_vars (
    id BIGSERIAL PRIMARY KEY,
    deployment_id BIGINT NOT NULL REFERENCES image_deployments(deployment_id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    value TEXT NOT NULL,

    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX image_environment_vars_deployment_id_name ON image_environment_vars(deployment_id, name);

CREATE TRIGGER trigger_image_environment_vars_updated_at
    BEFORE UPDATE ON image_environment_vars
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_image_environment_vars_updated_at ON image_environment_vars;
DROP TABLE image_environment_vars;