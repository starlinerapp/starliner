package application

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

const ingressDNSTimeout = 10 * time.Minute

type IngressApplication struct {
	deploy       port.Deploy
	dns          port.DNS
	queue        port.Queue
	logPublisher port.LogPublisher
}

func NewIngressApplication(
	deploy port.Deploy,
	dns port.DNS,
	queue port.Queue,
	logPublisher port.LogPublisher,
) *IngressApplication {
	return &IngressApplication{
		deploy:       deploy,
		dns:          dns,
		queue:        queue,
		logPublisher: logPublisher,
	}
}

func (ia *IngressApplication) HandleDeployIngress(i *value.IngressDeployment) {
	releaseName := i.DeploymentName
	hosts := toPortIngressHosts(i.IngressHosts)

	var logBuf strings.Builder
	appendStatus := ia.appendStatus(i.DeploymentId, i.Namespace, releaseName, &logBuf)

	appendStatus("==> Deploying ExternalDNS...\n")

	err := ia.deploy.DeployExternalDNS(i.Namespace, "external-dns", i.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to deploy external dns: %v\n", err)
	}

	if i.ExpectedIP == "" {
		appendStatus("==> ERROR: failed to deploy ExternalDNS: IP Address is not set\n")
		ia.publishDeploymentCompleted(i, logBuf.String())
		return
	}

	hostnames := ingressHostnames(i)
	args := newDeployIngressArgs(i, releaseName, hosts)

	if !ia.dns.AllHostsResolve(hostnames, i.ExpectedIP) {
		appendStatus("==> Deploying ingress without TLS for DNS propagation...\n")
		args.TLSEnabled = false
		if err := ia.deploy.DeployIngress(args); err != nil {
			appendStatus("==> ERROR: failed to deploy ingress: %v\n", err)
			ia.publishDeploymentCompleted(i, logBuf.String())
			return
		}
	}

	next := *i
	next.AccumulatedLogs = logBuf.String()
	if err := ia.queue.PublishEnableIngressTLS(&next); err != nil {
		log.Printf("failed to publish enable ingress tls: %v\n", err)
	}
}

func (ia *IngressApplication) HandleEnableIngressTLS(i *value.IngressDeployment) {
	deployment := *i
	go ia.enableIngressTLS(deployment)
}

func (ia *IngressApplication) enableIngressTLS(i value.IngressDeployment) {
	releaseName := i.DeploymentName
	hosts := toPortIngressHosts(i.IngressHosts)
	hostnames := ingressHostnames(&i)
	args := newDeployIngressArgs(&i, releaseName, hosts)

	var logBuf strings.Builder
	logBuf.WriteString(i.AccumulatedLogs)
	appendStatus := ia.appendStatus(i.DeploymentId, i.Namespace, releaseName, &logBuf)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("enable ingress TLS panic: %v", r)
		}
	}()

	if i.ExpectedIP == "" {
		appendStatus("==> ERROR: failed to deploy ExternalDNS: IP Address is not set\n")
		ia.publishDeploymentCompleted(&i, logBuf.String())
		return
	}

	if !ia.dns.AllHostsResolve(hostnames, i.ExpectedIP) {
		ctx, cancel := context.WithTimeout(context.Background(), ingressDNSTimeout)
		defer cancel()

		for _, host := range hostnames {
			appendStatus("==> Waiting for DNS propagation: %s -> %s\n", host, i.ExpectedIP)
			if err := ia.dns.WaitForHost(ctx, host, i.ExpectedIP); err != nil {
				appendStatus("==> ERROR: Failed to wait for DNS for %s: %v\n", host, err)
				retry := i
				retry.AccumulatedLogs = logBuf.String()
				if err := ia.queue.PublishEnableIngressTLS(&retry); err != nil {
					log.Printf("failed to republish enable ingress tls: %v\n", err)
				}
				return
			}
		}
	}

	appendStatus("==> Deploying ingress with TLS...\n")
	args.TLSEnabled = true
	if err := ia.deploy.DeployIngress(args); err != nil {
		appendStatus("==> ERROR: failed to deploy ingress: %v\n", err)
		ia.publishDeploymentCompleted(&i, logBuf.String())
		return
	}
	appendStatus("==> Ingress deployed successfully\n")
	ia.publishDeploymentCompleted(&i, logBuf.String())
}

func (ia *IngressApplication) appendStatus(
	deploymentId int64,
	namespace, releaseName string,
	logBuf *strings.Builder,
) func(format string, args ...any) {
	return func(format string, args ...any) {
		line := fmt.Sprintf(format, args...)
		logBuf.WriteString(line)
		if ia.logPublisher == nil {
			return
		}
		if err := ia.logPublisher.PublishLogChunk(context.Background(), deploymentId, namespace, releaseName, []byte(line)); err != nil {
			log.Printf("failed to publish log chunk: %v", err)
		}
	}
}

func (ia *IngressApplication) publishDeploymentCompleted(i *value.IngressDeployment, logs string) {
	if err := ia.queue.PublishDeploymentStatusLogsCompleted(&value.DeploymentStatusLogsCompleted{
		DeploymentId: i.DeploymentId,
		Logs:         logs,
	}); err != nil {
		log.Printf("failed to publish ingress deployment completed: %v", err)
	}
}

func ingressHostnames(i *value.IngressDeployment) []string {
	hostnames := make([]string, 0, len(i.IngressHosts))
	for _, host := range i.IngressHosts {
		hostnames = append(hostnames, host.Host)
	}
	return hostnames
}

func newDeployIngressArgs(
	i *value.IngressDeployment,
	releaseName string,
	hosts []port.IngressHost,
) *port.DeployIngressArgs {
	return &port.DeployIngressArgs{
		Namespace:        i.Namespace,
		ReleaseName:      releaseName,
		KubeconfigBase64: i.KubeconfigBase64,
		Hosts:            hosts,
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
