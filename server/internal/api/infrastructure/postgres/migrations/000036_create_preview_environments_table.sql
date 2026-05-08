-- +goose Up
CREATE TABLE preview_environments (
    environment_id bigint NOT NULL REFERENCES environments (id) ON DELETE CASCADE,
    github_repository_id bigint NOT NULL,
    pr_number bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT preview_environments_github_repository_id_pr_number_key UNIQUE (github_repository_id, pr_number)
);

CREATE TRIGGER trigger_preview_environments_updated_at
    BEFORE UPDATE ON preview_environments
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column ();

-- +goose Down
DROP TRIGGER trigger_preview_environments_updated_at ON preview_environments;

DROP TABLE preview_environments;

