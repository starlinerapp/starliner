-- +goose Up
CREATE TABLE users (
    id bigserial PRIMARY KEY,
    better_auth_id varchar(36) NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_better_auth_id ON users (better_auth_id);

CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;

DROP INDEX IF EXISTS idx_users_better_auth_id;

DROP TABLE users;

