-- +goose Up
CREATE TABLE image_deployments (
    deployment_id bigint PRIMARY KEY REFERENCES deployments (id) ON DELETE CASCADE,
    name varchar(255) NOT NULL,
    tag varchar(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_image_deployments_updated_at
    BEFORE UPDATE ON image_deployments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_image_deployments_updated_at ON image_deployments;

DROP TABLE image_deployments;

