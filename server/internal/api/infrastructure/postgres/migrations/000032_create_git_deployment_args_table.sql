-- +goose Up
CREATE TABLE git_deployment_args (
    id bigserial PRIMARY KEY,
    deployment_id bigint NOT NULL REFERENCES deployments (id) ON DELETE CASCADE,
    name text NOT NULL,
    value text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_git_deployment_args_updated_at
    BEFORE UPDATE ON git_deployment_args
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_git_deployment_args_updated_at ON git_deployment_args;

DROP TABLE git_deployment_args;

