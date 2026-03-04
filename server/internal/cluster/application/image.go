package application

import (
	"log"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type ImageApplication struct {
	deploy port.Deploy
}

func NewImageApplication(deploy port.Deploy) *ImageApplication {
	return &ImageApplication{deploy: deploy}
}

func (ia *ImageApplication) HandleDeployImage(a *value.ImageDeployment) {
	portEnvs := make([]*port.EnvVar, 0, len(a.EnvVars))
	for _, e := range a.EnvVars {
		portEnvs = append(portEnvs, &port.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		})
	}

	args := &port.DeployImageArgs{
		ReleaseName:      a.DeploymentName,
		KubeconfigBase64: a.KubeconfigBase64,
		Namespace:        a.Namespace,
		ImageRepository:  a.ImageName,
		ImageTag:         a.ImageTag,
		Port:             a.Port,
		EnvVars:          portEnvs,
	}

	err := ia.deploy.DeployImage(args)
	if err != nil {
		log.Printf("failed to deploy application: %v\n", err)
	}
}
