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
	hosts := toPortIngressHosts(i.IngressHosts)

	err := ia.deploy.DeployIngress(&port.DeployIngressArgs{
		ReleaseName:      releaseName,
		KubeconfigBase64: i.KubeconfigBase64,
		Hosts:            hosts,
	})
	if err != nil {
		log.Printf("failed to deploy ingress: %v\n", err)
	}
}

func toPortIngressHosts(in []value.IngressHost) []port.IngressHost {
	out := make([]port.IngressHost, 0, len(in))

	for _, h := range in {
		paths := make([]port.IngressPath, 0, len(h.Paths))
		for _, p := range h.Paths {
			paths = append(paths, port.IngressPath{
				Path:        p.Path,
				PathType:    port.PathType(p.PathType),
				ServiceName: p.ServiceName,
				ServicePort: p.ServicePort,
			})
		}

		out = append(out, port.IngressHost{
			Host:  h.Host,
			Paths: paths,
		})
	}

	return out
}
