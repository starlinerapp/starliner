-- +goose Up
ALTER TABLE deployments
    DROP CONSTRAINT deployments_environment_id_fkey;

ALTER TABLE deployments
    ALTER COLUMN environment_id DROP NOT NULL;

ALTER TABLE deployments
    ADD CONSTRAINT deployments_environment_id_fkey FOREIGN KEY (environment_id) REFERENCES environments (id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE deployments
    DROP CONSTRAINT deployments_environment_id_fkey;

ALTER TABLE deployments
    ALTER COLUMN environment_id SET NOT NULL;

ALTER TABLE deployments
    ADD CONSTRAINT deployments_environment_id_fkey FOREIGN KEY (environment_id) REFERENCES environments (id) ON DELETE RESTRICT;

