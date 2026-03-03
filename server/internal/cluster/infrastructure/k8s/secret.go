package k8s

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/infrastructure/shared/kubeconfig"
)

type Secret struct{}

func NewSecret() port.Secret {
	return &Secret{}
}

func (s *Secret) GetDatabaseCredentials(releaseName string, kubeconfigBase64 string) (*port.DatabaseCredentials, error) {
	var creds *port.DatabaseCredentials

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

		secret, err := client.CoreV1().
			Secrets(namespace).
			Get(ctx, releaseName+"-app", metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("get secret failed: %w", err)
		}

		username, ok := secret.Data["username"]
		if !ok {
			return fmt.Errorf("username not found in secret")
		}
		password, ok := secret.Data["password"]
		if !ok {
			return fmt.Errorf("password not found in secret")
		}
		dbName, ok := secret.Data["dbname"]
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
