-- +goose Up
CREATE TABLE runners (
  id BIGSERIAL PRIMARY KEY,
  organization_id BIGINT REFERENCES organizations (id) ON DELETE CASCADE,
  name TEXT,
  status TEXT NOT NULL DEFAULT 'offline',
  labels TEXT[] NOT NULL DEFAULT '{}',
  max_concurrent_jobs INTEGER NOT NULL DEFAULT 1,
  disabled_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT runners_max_concurrent_jobs_check CHECK (max_concurrent_jobs > 0)
);

CREATE TRIGGER trigger_runners_updated_at
  BEFORE UPDATE ON runners
  FOR EACH ROW
  EXECUTE PROCEDURE update_updated_at_column ();

CREATE TABLE runner_registration_tokens (
  id BIGSERIAL PRIMARY KEY,
  organization_id BIGINT REFERENCES organizations (id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,
  runner_id BIGINT REFERENCES runners (id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_runner_registration_tokens_updated_at
  BEFORE UPDATE ON runner_registration_tokens
  FOR EACH ROW
  EXECUTE PROCEDURE update_updated_at_column ();

-- +goose Down
DROP TABLE runner_registration_tokens;

DROP TABLE runners;
