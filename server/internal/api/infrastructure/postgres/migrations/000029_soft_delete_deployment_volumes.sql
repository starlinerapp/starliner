-- +goose Up
ALTER TABLE deployment_volumes ADD COLUMN deleted_at timestamptz;

ALTER TABLE deployment_volumes
    DROP CONSTRAINT deployment_volumes_deployment_id_fkey,
    ALTER COLUMN deployment_id DROP NOT NULL,
    ADD CONSTRAINT deployment_volumes_deployment_id_fkey
        FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE SET NULL;

-- +goose Down
DELETE FROM deployment_volumes WHERE deployment_id IS NULL;

ALTER TABLE deployment_volumes
    DROP CONSTRAINT deployment_volumes_deployment_id_fkey,
    ALTER COLUMN deployment_id SET NOT NULL,
    ADD CONSTRAINT deployment_volumes_deployment_id_fkey
        FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE;

ALTER TABLE deployment_volumes DROP COLUMN deleted_at;

