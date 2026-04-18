-- +goose Up
CREATE TABLE deployment_volumes (
    id BIGSERIAL PRIMARY KEY,
    deployment_id BIGINT REFERENCES deployments(id) ON DELETE SET NULL,
    volume_size_mib INTEGER NOT NULL,
    mount_path TEXT NOT NULL,
    deleted_at timestamptz,

    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_deployment_volumes_updated_at
    BEFORE UPDATE ON deployment_volumes
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_deployment_volumes_updated_at ON deployment_volumes;
DROP TABLE deployment_volumes;
