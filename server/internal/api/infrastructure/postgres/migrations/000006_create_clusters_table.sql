-- +goose Up
CREATE TYPE cluster_status AS ENUM (
    'pending',
    'running',
    'deleted'
);

CREATE TABLE clusters (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    ipv4_address varchar(255),
    public_key text,
    private_key text,
    organization_id bigint NOT NULL REFERENCES organizations (id),
    provisioning_id text,
    status cluster_status NOT NULL DEFAULT 'pending',
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_clusters_updated_at
    BEFORE UPDATE ON clusters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_clusters_updated_at ON clusters;

DROP TABLE clusters;

DROP TYPE IF EXISTS cluster_status;

