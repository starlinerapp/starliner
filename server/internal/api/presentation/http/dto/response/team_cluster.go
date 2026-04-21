package response

import "starliner.app/internal/api/domain/value"

type TeamCluster struct {
	TeamID      int64  `json:"teamId" binding:"required"`
	ClusterID   int64  `json:"clusterId" binding:"required"`
	ClusterName string `json:"clusterName" binding:"required"`
	ServerType  string `json:"serverType" binding:"required"`
}

func NewTeamCluster(tc *value.TeamCluster) TeamCluster {
	return TeamCluster{
		TeamID:      tc.TeamId,
		ClusterID:   tc.ClusterId,
		ClusterName: tc.ClusterName,
		ServerType:  tc.ServerType,
	}
}

func NewTeamClusters(tcs []*value.TeamCluster) []TeamCluster {
	clusters := make([]TeamCluster, len(tcs))
	for i, tc := range tcs {
		clusters[i] = NewTeamCluster(tc)
	}
	return clusters
}
