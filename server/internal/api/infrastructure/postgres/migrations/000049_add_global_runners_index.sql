-- +goose Up
CREATE INDEX idx_runners_global ON runners (id)
WHERE organization_id IS NULL;

-- +goose Down
DROP INDEX idx_runners_global;
