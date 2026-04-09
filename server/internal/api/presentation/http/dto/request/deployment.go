package request

type EnvVar struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type DeployFromGit struct {
	EnvironmentId         int64    `json:"environmentId" binding:"required"`
	ServiceName           string   `json:"serviceName" binding:"required"`
	Port                  int      `json:"port" binding:"required"`
	GitUrl                string   `json:"gitUrl" binding:"required"`
	ProjectRepositoryPath string   `json:"projectRepositoryPath" binding:"required"`
	DockerfilePath        string   `json:"dockerfilePath" binding:"required"`
	Envs                  []EnvVar `json:"envs" binding:"required"`
}

type UpdateDeployFromGit struct {
	EnvironmentId         int64    `json:"environmentId" binding:"required"`
	Port                  int      `json:"port" binding:"required"`
	ProjectRepositoryPath string   `json:"projectRepositoryPath" binding:"required"`
	DockerfilePath        string   `json:"dockerfilePath" binding:"required"`
	Envs                  []EnvVar `json:"envs" binding:"required"`
}

type DeployImage struct {
	EnvironmentId   int64    `json:"environmentId" binding:"required"`
	ServiceName     string   `json:"serviceName" binding:"required"`
	ImageName       string   `json:"imageName" binding:"required"`
	Tag             string   `json:"tag" binding:"required"`
	Port            int      `json:"port" binding:"required"`
	VolumeSizeMiB   *int32   `json:"volumeSizeMiB"`
	VolumeMountPath *string  `json:"volumeMountPath"`
	Envs            []EnvVar `json:"envs" binding:"required"`
}

type UpdateImage struct {
	EnvironmentId int64    `json:"environmentId" binding:"required"`
	ImageName     string   `json:"imageName" binding:"required"`
	Tag           string   `json:"tag" binding:"required"`
	Port          int      `json:"port" binding:"required"`
	Envs          []EnvVar `json:"envs" binding:"required"`
}

type DeployDatabase struct {
	EnvironmentId int64  `json:"environmentId" binding:"required"`
	ServiceName   string `json:"serviceName" binding:"required"`
}

type IngressPath struct {
	Path        string `json:"path" binding:"required"`
	PathType    string `json:"pathType" binding:"required,oneof=Prefix Exact"`
	ServiceName string `json:"serviceName" binding:"required"`
}

type IngressHost struct {
	Host  string        `json:"host" binding:"required"`
	Paths []IngressPath `json:"paths" binding:"required"`
}

type DeployIngress struct {
	EnvironmentId int64         `json:"environmentId" binding:"required"`
	IngressHosts  []IngressHost `json:"ingressHosts" binding:"required"`
}

type UpdateIngress struct {
	EnvironmentId int64         `json:"environmentId" binding:"required"`
	IngressHosts  []IngressHost `json:"ingressHosts" binding:"required"`
}
