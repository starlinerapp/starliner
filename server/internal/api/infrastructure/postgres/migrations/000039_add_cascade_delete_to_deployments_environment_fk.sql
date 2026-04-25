-- +goose Up
ALTER TABLE deployments
    DROP CONSTRAINT IF EXISTS deployments_environment_id_fkey;

ALTER TABLE deployments
    ADD CONSTRAINT deployments_environment_id_fkey
        FOREIGN KEY (environment_id)
            REFERENCES environments(id)
            ON DELETE CASCADE;

-- +goose Down
ALTER TABLE deployments
    DROP CONSTRAINT IF EXISTS deployments_environment_id_fkey;

ALTER TABLE deployments
    ADD CONSTRAINT deployments_environment_id_fkey
        FOREIGN KEY (environment_id)
            REFERENCES environments(id)
            ON DELETE RESTRICT;