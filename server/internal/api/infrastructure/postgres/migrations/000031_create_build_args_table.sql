-- +goose Up
CREATE TABLE build_args (
    id BIGSERIAL PRIMARY KEY,
    build_id BIGINT NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    value TEXT NOT NULL,

    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX build_args_build_id_name ON build_args(build_id, name);

CREATE TRIGGER trigger_build_args_updated_at
    BEFORE UPDATE ON build_args
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_build_args_updated_at ON build_args;
DROP TABLE build_args;

