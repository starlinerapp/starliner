package request

type TeamClusterAssignment struct {
	ClusterId int64 `json:"clusterId" binding:"required"`
}

type SetTeamClusters struct {
	Clusters []TeamClusterAssignment `json:"clusters"`
}
