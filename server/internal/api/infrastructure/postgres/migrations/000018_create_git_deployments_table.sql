-- +goose Up
CREATE TABLE git_deployments (
    deployment_id BIGINT PRIMARY KEY REFERENCES deployments(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    project_path TEXT NOT NULL,
    dockerfile_path TEXT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_git_deployments_updated_at
    BEFORE UPDATE ON git_deployments
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_git_deployments_updated_at ON git_deployments;
DROP TABLE IF EXISTS git_deployments;