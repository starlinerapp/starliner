-- +goose Up
ALTER TABLE organization_invites
    ADD COLUMN email text NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE organization_invites
    DROP COLUMN email;

