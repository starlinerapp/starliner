package application

import (
	"fmt"
	"log"

	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type IngressApplication struct {
	deploy   port.Deploy
	queue    port.Queue
	notifier *Notifier
}

func NewIngressApplication(deploy port.Deploy, queue port.Queue) *IngressApplication {
	notifier := NewNotifier(queue)
	return &IngressApplication{
		deploy:   deploy,
		queue:    queue,
		notifier: notifier,
	}
}

func (ia *IngressApplication) HandleDeployIngress(i *value.IngressDeployment) {

	if i.CorrelationId == nil {
		log.Printf("missing correlation id for ingress deployment %d\n", i.DeploymentId)
	}

	releaseName := i.DeploymentName
	hosts := toPortIngressHosts(i.IngressHosts)

	err := ia.deploy.DeployIngress(&port.DeployIngressArgs{
		Namespace:        i.Namespace,
		ReleaseName:      releaseName,
		KubeconfigBase64: i.KubeconfigBase64,
		Hosts:            hosts,
	})
	if err != nil {
		log.Printf("failed to deploy ingress: %v\n", err)
		ia.notifier.publishNotification(i.DeploymentId, *i.CorrelationId, "failed", fmt.Sprintf("Failed to deploy ingress: %s", i.DeploymentName))
		return
	}
	ia.notifier.publishNotification(i.DeploymentId, *i.CorrelationId, "success", fmt.Sprintf("Successfully deployed ingress: %s", i.DeploymentName))
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
