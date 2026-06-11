package value

type EnvVar struct {
	Name  string
	Value string
}

type ImageDeployment struct {
	DeploymentId          int64
	CorrelationId         *string
	DeploymentName        string
	Namespace             string
	KubeconfigBase64      string
	ImageRegistryUrl      string
	ImageRegistryUsername string
	ImageRegistryPassword string
	ImageName             string
	ImageTag              string
	Port                  int
	VolumeSizeMiB         *int32
	VolumeMountPath       *string
	EnvVars               []*EnvVar
}

type Deployment struct {
	Namespace        string
	DeploymentId     int64
	CorrelationId    *string
	DeploymentName   string
	KubeconfigBase64 string
	ClusterId        int64
	OrganizationId   int64
	ProvisioningId   string
}

type PathType string

const (
	Prefix PathType = "Prefix"
	Exact  PathType = "Exact"
)

type IngressPath struct {
	Path        string
	PathType    PathType
	ServiceName string
	ServicePort int
}

type IngressHost struct {
	Host  string
	Paths []IngressPath
}

type IngressDeployment struct {
	DeploymentId     int64
	CorrelationId    *string
	DeploymentName   string
	Namespace        string
	KubeconfigBase64 string
	ExpectedIP       string
	IngressHosts     []IngressHost
	AccumulatedLogs  string
}

type DeploymentStatusLogsCompleted struct {
	DeploymentId int64
	Logs         string
}

type DatabaseDeployedSuccess struct {
	CorrelationId  string `json:"correlationId"`
	DeploymentId   int64  `json:"deploymentId"`
	DeploymentName string `json:"deploymentName"`
	DbName         string `json:"dbName"`
	Username       string `json:"username"`
	Password       string `json:"password"`
}

type DatabaseDeployedFailure struct {
	CorrelationId  string `json:"correlationId"`
	DeploymentId   int64  `json:"deploymentId"`
	DeploymentName string `json:"deploymentName"`
}

type ImageDeployedSuccess struct {
	CorrelationId string `json:"correlationId"`
	DeploymentId  int64  `json:"deploymentId"`
	ImageName     string `json:"imageName"`
}

type ImageDeployedFailure struct {
	CorrelationId string `json:"correlationId"`
	DeploymentId  int64  `json:"deploymentId"`
	ImageName     string `json:"imageName"`
}

type IngressDeployedSuccess struct {
	CorrelationId  string `json:"correlationId"`
	DeploymentId   int64  `json:"deploymentId"`
	DeploymentName string `json:"deploymentName"`
}

type IngressDeployedFailure struct {
	CorrelationId  string `json:"correlationId"`
	DeploymentId   int64  `json:"deploymentId"`
	DeploymentName string `json:"deploymentName"`
}

type DeploymentDeletedSuccess struct {
	CorrelationId  string `json:"correlationId"`
	DeploymentId   int64  `json:"deploymentId"`
	DeploymentName string `json:"deploymentName"`
}

type DeploymentDeletedFailure struct {
	CorrelationId  string `json:"correlationId"`
	DeploymentId   int64  `json:"deploymentId"`
	DeploymentName string `json:"deploymentName"`
}
