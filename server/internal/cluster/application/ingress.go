package application

import (
	"log"
	"starliner.app/internal/cluster/domain/port"
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
	releaseName := i.DeploymentName
	err := ia.deploy.DeployIngress(&port.DeployIngressArgs{
		ReleaseName:      releaseName,
		KubeconfigBase64: i.KubeconfigBase64,
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
	if err != nil {
		log.Printf("failed to deploy ingress: %v\n", err)
	}
}
