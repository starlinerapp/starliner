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
	releaseName := a.DeploymentName
	err := ia.deploy.DeployImage(
		releaseName,
		a.KubeconfigBase64,
		a.ImageRepository,
		a.ImageTag,
		a.Port,
	)
	if err != nil {
		log.Printf("failed to deploy application: %v\n", err)
	}
}
