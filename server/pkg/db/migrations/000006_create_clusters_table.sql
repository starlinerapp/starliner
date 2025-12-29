-- +goose Up
CREATE TABLE clusters (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    ipv4_address VARCHAR(255) NOT NULL,
    public_key TEXT not null,
    private_key_ref VARCHAR(255) NOT NULL,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_clusters_updated_at
    BEFORE UPDATE ON clusters
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_clusters_updated_at ON clusters;

DROP TABLE clusters;