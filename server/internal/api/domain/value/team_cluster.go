package value

import "starliner.app/internal/api/domain/entity"

type TeamCluster struct {
	TeamId      int64
	ClusterId   int64
	ClusterName string
	ServerType  string
}

func NewTeamCluster(tc *entity.TeamCluster) *TeamCluster {
	return &TeamCluster{
		TeamId:      tc.TeamId,
		ClusterId:   tc.ClusterId,
		ClusterName: tc.ClusterName,
		ServerType:  tc.ServerType,
	}
}

func NewTeamClusters(tcs []*entity.TeamCluster) []*TeamCluster {
	clusters := make([]*TeamCluster, len(tcs))
	for i, cluster := range tcs {
		clusters[i] = NewTeamCluster(cluster)
	}
	return clusters
}
