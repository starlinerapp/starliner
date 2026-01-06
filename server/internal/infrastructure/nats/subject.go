package nats

type Subject string

const (
	BuildTriggered Subject = "build.triggered"
	CreateCluster  Subject = "create.cluster"
	DeleteCluster  Subject = "delete.cluster"
	CreateProject  Subject = "create.project"
)
