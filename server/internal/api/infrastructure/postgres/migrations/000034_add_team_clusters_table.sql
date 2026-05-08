-- +goose Up
CREATE TABLE team_clusters (
    team_id bigint NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    cluster_id bigint NOT NULL REFERENCES clusters (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (team_id, cluster_id)
);

CREATE TRIGGER trigger_team_clusters_updated_at
    BEFORE UPDATE ON team_clusters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column ();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_team_clusters_updated_at ON team_clusters;

DROP TABLE team_clusters;

