package application

import (
	"log"

	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type ImageApplication struct {
	deploy port.Deploy
	queue  port.Queue
}

func NewImageApplication(deploy port.Deploy, queue port.Queue) *ImageApplication {
	return &ImageApplication{deploy: deploy, queue: queue}
}

func (ia *ImageApplication) HandleDeployImage(a *value.ImageDeployment) {
	correlationId := ""
	if a.CorrelationId != nil {
		correlationId = *a.CorrelationId
	} else {
		log.Printf("missing correlation id for image deployment %d\n", a.DeploymentId)
	}

	portEnvs := make([]*port.EnvVar, 0, len(a.EnvVars))
	for _, e := range a.EnvVars {
		portEnvs = append(portEnvs, &port.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		})
	}

	args := &port.DeployImageArgs{
		ReleaseName:           a.DeploymentName,
		KubeconfigBase64:      a.KubeconfigBase64,
		Namespace:             a.Namespace,
		ImageName:             a.ImageName,
		ImageRegistryUrl:      a.ImageRegistryUrl,
		ImageRegistryUsername: a.ImageRegistryUsername,
		ImageRegistryPassword: a.ImageRegistryPassword,
		ImageTag:              a.ImageTag,
		Port:                  a.Port,
		VolumeSizeMiB:         a.VolumeSizeMiB,
		VolumeMountPath:       a.VolumeMountPath,
		EnvVars:               portEnvs,
	}

	err := ia.deploy.DeployImage(args)
	if err != nil {
		log.Printf("failed to deploy application: %v\n", err)
		if pubErr := ia.queue.PublishImageDeployedFailure(&value.ImageDeployedFailure{
			CorrelationId: correlationId,
			DeploymentId:  a.DeploymentId,
			ImageName:     a.ImageName,
		}); pubErr != nil {
			log.Printf("failed to publish image deployed failure: %v\n", pubErr)
		}
		return
	}

	if pubErr := ia.queue.PublishImageDeployedSuccess(&value.ImageDeployedSuccess{
		CorrelationId: correlationId,
		DeploymentId:  a.DeploymentId,
		ImageName:     a.ImageName,
	}); pubErr != nil {
		log.Printf("failed to publish image deployed success: %v\n", pubErr)
	}
}
