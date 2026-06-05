package k8s

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/domain/value"
	"starliner.app/internal/cluster/infrastructure/shared/kubeconfig"
)

const (
	traefikNamespace     = "kube-system"
	traefikLabelSelector = "app.kubernetes.io/name=traefik"
)

type Logs struct{}

func NewLogs() port.Logs {
	return &Logs{}
}

func (l *Logs) StreamLogs(
	ctx context.Context,
	source value.LogSource,
	environmentNamespace string,
	releaseName string,
	kubeconfigBase64 string,
) (io.ReadCloser, error) {
	switch source {
	case value.LogSourceIngress:
		return l.streamIngressLogs(ctx, environmentNamespace, releaseName, kubeconfigBase64)
	default:
		return l.streamWorkloadLogs(ctx, environmentNamespace, releaseName, kubeconfigBase64)
	}
}

func (l *Logs) streamWorkloadLogs(
	ctx context.Context,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
) (io.ReadCloser, error) {
	client, err := newKubernetesClient(kubeconfigBase64)
	if err != nil {
		return nil, err
	}

	labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", releaseName)
	return streamPodLogs(ctx, client, namespace, labelSelector, nil)
}

func (l *Logs) streamIngressLogs(
	ctx context.Context,
	environmentNamespace string,
	releaseName string,
	kubeconfigBase64 string,
) (io.ReadCloser, error) {
	client, err := newKubernetesClient(kubeconfigBase64)
	if err != nil {
		return nil, err
	}

	filter := ingressLogLineFilter(environmentNamespace, releaseName)
	return streamPodLogs(ctx, client, traefikNamespace, traefikLabelSelector, filter)
}

func ingressLogLineFilter(environmentNamespace, releaseName string) func(string) bool {
	routerKey := fmt.Sprintf("%s-%s", environmentNamespace, releaseName)

	return func(line string) bool {
		return matchesTraefikRouterKey(line, routerKey)
	}
}

// matchesTraefikRouterKey matches Traefik access-log router names such as
// "{namespace}-{release}-{hostname-with-hyphens}@kubernetes".
func matchesTraefikRouterKey(line, routerKey string) bool {
	idx := strings.Index(line, routerKey)
	if idx < 0 {
		return false
	}

	after := idx + len(routerKey)
	if after >= len(line) {
		return true
	}

	switch line[after] {
	case '-', '@', '"', ' ':
		return true
	default:
		return false
	}
}

func newKubernetesClient(kubeconfigBase64 string) (*kubernetes.Clientset, error) {
	var client *kubernetes.Clientset

	err := kubeconfig.WithTempKubeConfig(kubeconfigBase64, func(kubeconfigPath string) error {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return fmt.Errorf("failed to build kubeconfig: %w", err)
		}

		client, err = kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func streamPodLogs(
	ctx context.Context,
	client *kubernetes.Clientset,
	namespace string,
	labelSelector string,
	lineFilter func(string) bool,
) (io.ReadCloser, error) {
	reader, writer := io.Pipe()

	go func() {
		defer func() {
			_ = writer.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				_ = writer.CloseWithError(ctx.Err())
				return
			default:
			}

			pods, err := client.CoreV1().
				Pods(namespace).
				List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				if ctx.Err() != nil {
					_ = writer.CloseWithError(ctx.Err())
					return
				}
				time.Sleep(2 * time.Second)
				continue
			}

			if len(pods.Items) == 0 {
				time.Sleep(2 * time.Second)
				continue
			}

			// TODO: merge logs of all matching pods
			pod := pods.Items[0]
			tailLines := int64(50)

			stream, err := client.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
				Follow:    true,
				TailLines: &tailLines,
			}).Stream(ctx)
			if err != nil {
				if ctx.Err() != nil {
					_ = writer.CloseWithError(ctx.Err())
					return
				}
				time.Sleep(2 * time.Second)
				continue
			}

			scanner := bufio.NewScanner(stream)
			for scanner.Scan() {
				line := scanner.Text()
				if lineFilter != nil && !lineFilter(line) {
					continue
				}

				if _, err := fmt.Fprintln(writer, line); err != nil {
					_ = stream.Close()
					return
				}
			}
			_ = stream.Close()

			if ctx.Err() != nil {
				_ = writer.CloseWithError(ctx.Err())
				return
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return reader, nil
}
