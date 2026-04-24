package k8s

import (
	"context"
	"errors"
	"fmt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/infrastructure/shared/kubeconfig"
	"time"
)

type Secret struct{}

func NewSecret() port.Secret {
	return &Secret{}
}

func (s *Secret) GetDatabaseCredentials(namespace string, releaseName string, kubeconfigBase64 string) (*port.DatabaseCredentials, error) {
	var creds *port.DatabaseCredentials

	err := kubeconfig.WithTempKubeConfig(kubeconfigBase64, func(kubeconfigPath string) error {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return fmt.Errorf("failed to build kubeconfig: %w", err)
		}

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		secretName := releaseName + "-db-credentials"

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		var secretData map[string][]byte

		pollErr := wait.ExponentialBackoffWithContext(ctx, wait.Backoff{
			Duration: 500 * time.Millisecond,
			Factor:   1.7,
			Jitter:   0.1,
			Steps:    20,
		}, func(ctx context.Context) (bool, error) {
			sec, err := client.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					return false, nil
				}
				return false, err
			}
			secretData = sec.Data
			return true, nil
		})

		if pollErr != nil {
			if errors.Is(pollErr, context.DeadlineExceeded) {
				return fmt.Errorf("timed out waiting for secret %s/%s", namespace, secretName)
			}
			return fmt.Errorf("failed waiting for secret: %w", pollErr)
		}

		username, ok := secretData["username"]
		if !ok {
			return fmt.Errorf("username not found in secret")
		}
		password, ok := secretData["password"]
		if !ok {
			return fmt.Errorf("password not found in secret")
		}
		dbName, ok := secretData["dbname"]
		if !ok {
			return fmt.Errorf("dbname not found in secret")
		}

		creds = &port.DatabaseCredentials{
			DatabaseName: string(dbName),
			Username:     string(username),
			Password:     string(password),
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return creds, nil
}
