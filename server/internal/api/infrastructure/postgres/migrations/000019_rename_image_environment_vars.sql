-- +goose Up
ALTER TABLE image_environment_vars RENAME TO deployment_environment_vars;

ALTER INDEX image_environment_vars_deployment_id_name RENAME TO deployment_environment_vars_deployment_id_name;

ALTER TRIGGER trigger_image_environment_vars_updated_at ON deployment_environment_vars RENAME TO trigger_deployment_environment_vars_updated_at;

ALTER TABLE deployment_environment_vars
  DROP CONSTRAINT image_environment_vars_deployment_id_fkey;

ALTER TABLE deployment_environment_vars
  ADD CONSTRAINT deployment_environment_vars_deployment_id_fkey FOREIGN KEY (deployment_id) REFERENCES deployments (id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE deployment_environment_vars
  DROP CONSTRAINT deployment_environment_vars_deployment_id_fkey;

ALTER TABLE deployment_environment_vars
  ADD CONSTRAINT image_environment_vars_deployment_id_fkey FOREIGN KEY (deployment_id) REFERENCES image_deployments (deployment_id) ON DELETE CASCADE;

ALTER TRIGGER trigger_deployment_environment_vars_updated_at ON deployment_environment_vars RENAME TO trigger_image_environment_vars_updated_at;

ALTER INDEX deployment_environment_vars_deployment_id_name RENAME TO image_environment_vars_deployment_id_name;

ALTER TABLE deployment_environment_vars RENAME TO image_environment_vars;

