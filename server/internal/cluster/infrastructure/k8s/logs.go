package k8s

import (
	"bufio"
	"context"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/infrastructure/shared/kubeconfig"
	"time"
)

type Logs struct{}

func NewLogs() port.Logs {
	return &Logs{}
}

func (l *Logs) StreamLogs(ctx context.Context, namespace string, releaseName string, kubeconfigBase64 string) (io.ReadCloser, error) {
	reader, writer := io.Pipe()
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

		go func() {
			defer func() {
				_ = writer.Close()
			}()

			labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", releaseName)

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

				stream, err := client.CoreV1().
					Pods(namespace).
					GetLogs(pod.Name, &corev1.PodLogOptions{
						Follow:    true,
						TailLines: &tailLines,
					}).
					Stream(ctx)
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
					_, err := fmt.Fprintln(writer, scanner.Text())
					if err != nil {
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
		return nil
	})
	if err != nil {
		return nil, err
	}

	return reader, nil
}
