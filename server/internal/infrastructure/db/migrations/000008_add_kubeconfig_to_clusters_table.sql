-- +goose Up
ALTER TABLE clusters
    ADD COLUMN kubeconfig TEXT;

-- +goose Down
ALTER TABLE clusters
    DROP COLUMN kubeconfig;