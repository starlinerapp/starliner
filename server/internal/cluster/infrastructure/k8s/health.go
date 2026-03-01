package k8s

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/domain/value"
	"starliner.app/internal/cluster/infrastructure/shared/kubeconfig"
)

type Health struct {
}

func NewHealth() port.Health {
	return &Health{}
}

func (h *Health) CheckPodsHealthy(releaseName string, kubeconfigBase64 string) (*value.HealthStatus, error) {
	var status *value.HealthStatus

	err := kubeconfig.WithTempKubeConfig(kubeconfigBase64, func(kubeconfigPath string) error {
		ctx := context.Background()

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return fmt.Errorf("failed to build kubeconfig: %w", err)
		}

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		const namespace = "default"
		labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", releaseName)

		pods, err := client.CoreV1().
			Pods(namespace).
			List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			return fmt.Errorf("list pods failed: %w", err)
		}

		if len(pods.Items) == 0 {
			status = &value.HealthStatus{
				Health: value.Unhealthy,
				Status: fmt.Sprintf(
					"no pods found for release %q in namespace %q",
					releaseName,
					namespace,
				),
			}
			return nil
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase != corev1.PodRunning {
				status = &value.HealthStatus{
					Health: value.Unhealthy,
					Status: fmt.Sprintf(
						"pod %s not running (phase=%s)",
						pod.Name,
						pod.Status.Phase,
					),
				}
				return nil
			}

			for _, cs := range pod.Status.ContainerStatuses {
				if !cs.Ready {
					status = &value.HealthStatus{
						Health: value.Unhealthy,
						Status: fmt.Sprintf(
							"pod %s container %s not ready",
							pod.Name,
							cs.Name,
						),
					}
					return nil
				}

				if cs.State.Waiting != nil {
					switch cs.State.Waiting.Reason {
					case "CrashLoopBackOff", "ImagePullBackOff", "ErrImagePull", "Error":
						reason := cs.State.Waiting.Reason
						status = &value.HealthStatus{
							Health: value.Unhealthy,
							Status: fmt.Sprintf(
								"pod %s container %s unhealthy (reason=%s)",
								pod.Name,
								cs.Name,
								reason,
							),
						}
						return nil
					}
				}
			}
		}

		status = &value.HealthStatus{
			Health: value.Healthy,
			Status: "all pods healthy",
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return status, nil
}
