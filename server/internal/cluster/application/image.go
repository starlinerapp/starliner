package application

import (
	"fmt"
	"log"

	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type ImageApplication struct {
	deploy   port.Deploy
	queue    port.Queue
	notifier *Notifier
}

func NewImageApplication(deploy port.Deploy, queue port.Queue) *ImageApplication {
	notifier := NewNotifier(queue)
	return &ImageApplication{deploy: deploy, queue: queue, notifier: notifier}
}

func (ia *ImageApplication) HandleDeployImage(a *value.ImageDeployment) {
	if a.CorrelationId == nil {
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
		ia.notifier.publishNotification(a.DeploymentId, *a.CorrelationId, "failed", fmt.Sprintf("Failed to deploy image: %s", a.ImageName))
		return
	}
	ia.notifier.publishNotification(a.DeploymentId, *a.CorrelationId, "success", fmt.Sprintf("Successfully deployed image: %s", a.ImageName))
}
