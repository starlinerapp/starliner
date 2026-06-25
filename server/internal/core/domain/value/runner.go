package value

type RunnerStatus string

const (
	RunnerStatusOnline  RunnerStatus = "online"
	RunnerStatusOffline RunnerStatus = "offline"
)

type RunnerStatusChanged struct {
	RunnerId int64        `json:"runner_id"`
	Status   RunnerStatus `json:"status"`
}

type RunnerDeleted struct {
	RunnerId int64 `json:"runner_id"`
}
