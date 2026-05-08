-- +goose Up
ALTER TABLE clusters
    ADD COLUMN kubeconfig text;

-- +goose Down
ALTER TABLE clusters
    DROP COLUMN kubeconfig;

