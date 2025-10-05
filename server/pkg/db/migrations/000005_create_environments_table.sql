-- +goose Up
CREATE TABLE environments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_environments_updated_at
    BEFORE UPDATE ON environments
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_environments_updated_at ON environments;

DROP TABLE environments;