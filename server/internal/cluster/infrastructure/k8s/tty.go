package k8s

import (
	"context"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/infrastructure/shared/kubeconfig"
	"time"
)

type TTY struct{}

func NewTTY() port.TTY {
	return &TTY{}
}

func (t *TTY) Open(
	ctx context.Context,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
	terminalSize <-chan port.TerminalSize,
) (stdin io.WriteCloser, stdout io.ReadCloser, err error) {
	var client *kubernetes.Clientset
	var restConfig *rest.Config

	err = kubeconfig.WithTempKubeConfig(kubeconfigBase64, func(kubeconfigPath string) error {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return fmt.Errorf("failed to build kubeconfig: %w", err)
		}

		restConfig = config
		client, err = kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	go func() {
		defer func() {
			_ = stdinReader.Close()
			_ = stdoutWriter.Close()
		}()

		labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", releaseName)

		var pod *corev1.Pod
		for {
			select {
			case <-ctx.Done():
				_ = stdoutWriter.CloseWithError(ctx.Err())
				return
			default:
			}

			pods, err := client.CoreV1().
				Pods(namespace).
				List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				if ctx.Err() != nil {
					_ = stdoutWriter.CloseWithError(ctx.Err())
				}
				time.Sleep(2 * time.Second)
				continue
			}

			if len(pods.Items) == 0 {
				time.Sleep(2 * time.Second)
				continue
			}

			pod = &pods.Items[0]
			break
		}

		req := client.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(pod.Name).
			Namespace(namespace).
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Stdin:   true,
				Stdout:  true,
				Stderr:  true,
				TTY:     true,
				Command: []string{"/bin/sh"},
			}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
		if err != nil {
			_ = stdoutWriter.CloseWithError(err)
			return
		}

		err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdin:             stdinReader,
			Stdout:            stdoutWriter,
			Stderr:            stdoutWriter,
			Tty:               true,
			TerminalSizeQueue: &terminalSizeQueue{sizes: terminalSize},
		})
		if err != nil {
			_ = stdoutWriter.CloseWithError(err)
			return
		}
	}()

	return stdinWriter, stdoutReader, nil
}

type terminalSizeQueue struct {
	sizes <-chan port.TerminalSize
}

func (q *terminalSizeQueue) Next() *remotecommand.TerminalSize {
	size, ok := <-q.sizes
	if !ok {
		return nil
	}
	return &remotecommand.TerminalSize{
		Width:  size.Columns,
		Height: size.Rows,
	}
}
