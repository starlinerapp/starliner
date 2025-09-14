CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    better_auth_id BIGSERIAL NOT NULL
);

CREATE INDEX ON users (better_auth_id)