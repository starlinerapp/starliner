package application

import (
	"log"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type ApplicationApplication struct {
	deploy port.Deploy
}

func NewApplicationApplication(deploy port.Deploy) *ApplicationApplication {
	return &ApplicationApplication{deploy: deploy}
}

func (aa *ApplicationApplication) HandleDeployApplication(a *value.ApplicationDeployment) {
	releaseName := a.DeploymentName
	err := aa.deploy.DeployApplication(
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
