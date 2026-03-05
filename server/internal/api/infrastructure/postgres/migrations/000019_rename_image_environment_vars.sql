-- +goose Up
ALTER TABLE image_environment_vars RENAME TO deployment_environment_vars;

-- +goose Down
ALTER TABLE deployment_environment_vars RENAME TO image_environment_vars;