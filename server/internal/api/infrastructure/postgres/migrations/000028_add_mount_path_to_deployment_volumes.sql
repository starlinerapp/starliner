-- +goose Up
ALTER TABLE deployment_volumes ADD COLUMN mount_path TEXT NOT NULL DEFAULT '/data';

-- +goose Down
ALTER TABLE deployment_volumes DROP COLUMN mount_path;

