-- +goose Up
CREATE TABLE github_apps (
    id bigserial PRIMARY KEY,
    installation_id bigint NOT NULL,
    organization_id bigint NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT github_apps_organization_id_key UNIQUE (organization_id)
);

CREATE TRIGGER trigger_github_apps_updated_at
    BEFORE UPDATE ON github_apps
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_github_apps_updated_at ON github_apps;

DROP TABLE github_apps;

