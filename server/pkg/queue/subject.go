package queue

type Subject string

const (
	BuildTriggered Subject = "build.triggered"
	CreateCluster  Subject = "create.cluster"
)
