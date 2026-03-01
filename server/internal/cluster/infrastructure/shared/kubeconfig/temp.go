package kubeconfig

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

func WithTempKubeConfig(kubeconfigBase64 string, fn func(path string) error) error {
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(kubeconfigBase64)
	if err != nil {
		return fmt.Errorf("failed to decode kubeconfig: %v", err)

	}

	tmpDir, err := os.MkdirTemp("", "kubeconfig-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("failed to remove temp directory: %v", err)
		}
	}()

	kubeconfigPath := filepath.Join(tmpDir, "kubeconfig")
	err = os.WriteFile(kubeconfigPath, kubeconfigBytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to write kubeconfig: %v", err)
	}

	return fn(kubeconfigPath)
}
