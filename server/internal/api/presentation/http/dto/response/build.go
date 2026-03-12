package response

type BuildLogs struct {
	Logs *string `json:"logs" binding:"required"`
}

func NewBuildLogs(logs *string) BuildLogs {
	return BuildLogs{
		Logs: logs,
	}
}
