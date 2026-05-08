-- +goose Up
CREATE TABLE projects (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    organization_id bigint NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (name, organization_id)
);

CREATE TRIGGER trigger_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_projects_updated_at ON projects;

DROP TABLE projects;

