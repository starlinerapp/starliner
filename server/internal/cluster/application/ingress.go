package application

import (
	"log"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/utils"
	"starliner.app/internal/core/domain/value"
)

type IngressApplication struct {
	deploy port.Deploy
}

func NewIngressApplication(deploy port.Deploy) *IngressApplication {
	return &IngressApplication{
		deploy: deploy,
	}
}

func (ia *IngressApplication) HandleDeployIngress(i *value.IngressDeployment) {
	err := utils.WithTempKubeConfig(i.KubeconfigBase64, func(kubeconfigPath string) error {
		releaseName := i.DeploymentName
		return ia.deploy.DeployIngress(&port.DeployIngressArgs{
			ReleaseName:    releaseName,
			KubeconfigPath: kubeconfigPath,
			Hosts: []port.IngressHost{
				{
					Host: "46.225.14.208.nip.io",
					Paths: []port.IngressPath{
						{
							Path:        "/",
							PathType:    "Prefix",
							ServiceName: "cnpg-webhook-service",
							ServicePort: 80,
						},
					},
				},
			},
		})
	})
	if err != nil {
		log.Printf("failed to deploy ingress: %v\n", err)
	}
}
