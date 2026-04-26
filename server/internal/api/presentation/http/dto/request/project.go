package request

type CreateProject struct {
	Name      string `json:"name" binding:"required"`
	TeamId    int64  `json:"team_id" binding:"required"`
	ClusterId int64  `json:"cluster_id" binding:"required"`
}
