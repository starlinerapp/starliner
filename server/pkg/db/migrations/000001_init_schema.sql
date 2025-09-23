-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    better_auth_id VARCHAR(36) NOT NULL
);

CREATE INDEX idx_users_better_auth_id ON users (better_auth_id);

-- +goose Down
DROP INDEX IF EXISTS idx_users_better_auth_id;

DROP TABLE users;
