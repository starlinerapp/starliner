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
)

type Health struct {
}

func NewHealth() port.Health {
	return &Health{}
}

func (h *Health) CheckPodsHealthy(releaseName string, kubeconfigPath string) (*value.HealthStatus, error) {
	ctx := context.Background()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", releaseName)

	pods, err := client.CoreV1().
		Pods("default").
		List(ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})
	if err != nil {
		return nil, fmt.Errorf("list pods failed: %w", err)
	}

	if len(pods.Items) == 0 {
		return &value.HealthStatus{
			Health: value.Unhealthy,
			Status: fmt.Sprintf(
				"no pods found for release %q in namespace %q",
				releaseName,
				"default",
			),
		}, nil
	}

	for _, pod := range pods.Items {
		// Pod phase must be Running
		if pod.Status.Phase != corev1.PodRunning {
			return &value.HealthStatus{
				Health: value.Unhealthy,
				Status: fmt.Sprintf(
					"pod %s not running (phase=%s)",
					pod.Name,
					pod.Status.Phase,
				),
			}, nil
		}

		for _, cs := range pod.Status.ContainerStatuses {
			// All containers must be Ready
			if !cs.Ready {
				return &value.HealthStatus{
					Health: value.Unhealthy,
					Status: fmt.Sprintf(
						"pod %s container %s not ready",
						pod.Name,
						cs.Name,
					),
				}, nil
			}

			// No bad waiting states
			if cs.State.Waiting != nil {
				reason := cs.State.Waiting.Reason
				switch reason {
				case "CrashLoopBackOff",
					"ImagePullBackOff",
					"ErrImagePull",
					"Error":
					return &value.HealthStatus{
						Health: value.Unhealthy,
						Status: fmt.Sprintf(
							"pod %s container %s unhealthy (reason=%s)",
							pod.Name,
							cs.Name,
							reason,
						),
					}, nil
				}
			}
		}
	}
	return &value.HealthStatus{
		Health: value.Healthy,
		Status: "all pods healthy",
	}, nil
}
