package value

type EnvironmentNotification struct {
	CorrelationId string `json:"correlationId"`
	DeploymentId  int64  `json:"deploymentId"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type ClusterNotification struct {
	ClusterId int64  `json:"clusterId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
