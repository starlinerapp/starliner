package request

type RegisterRunner struct {
	Token             string   `json:"token" binding:"required"`
	Name              string   `json:"name" binding:"required"`
	Labels            []string `json:"labels"`
	MaxConcurrentJobs int32    `json:"maxConcurrentJobs" binding:"required,gt=0"`
}
