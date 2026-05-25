package application

import (
	"context"
	"log"
	"time"

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

	log.Printf("deploying external DNS")
	err := ia.deploy.DeployExternalDNS(i.Namespace, "external-dns", i.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to deploy external dns: %v\n", err)
	}

	args := &port.DeployIngressArgs{
		Namespace:        i.Namespace,
		ReleaseName:      releaseName,
		KubeconfigBase64: i.KubeconfigBase64,
		Hosts:            hosts,
	}

	hostnames := make([]string, 0, len(i.IngressHosts))
	for _, host := range i.IngressHosts {
		hostnames = append(hostnames, host.Host)
	}

	if i.ExpectedIP == "" {
		log.Printf("failed to deploy ingress: cluster expected IP is not set\n")
		return
	}

	if !allHostsResolve(hostnames, i.ExpectedIP) {
		log.Printf("deploying ingress without TLS for DNS propagation")
		args.TLSEnabled = false
		if err := ia.deploy.DeployIngress(args); err != nil {
			log.Printf("failed to deploy ingress (phase 1): %v\n", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		for _, host := range hostnames {
			log.Printf("waiting for DNS propagation: %s -> %s", host, i.ExpectedIP)
			if err := waitForDNS(ctx, host, i.ExpectedIP); err != nil {
				log.Printf("failed to wait for DNS for %s: %v\n", host, err)
				return
			}
		}
	}

	log.Printf("deploying ingress with TLS")
	args.TLSEnabled = true
	if err := ia.deploy.DeployIngress(args); err != nil {
		log.Printf("failed to deploy ingress (phase 2): %v\n", err)
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
