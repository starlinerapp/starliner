package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

// TODO: Move this logic to infrastructure

func WithTempKubeConfig(kubeconfigBase64 string, fn func(path string) error) error {
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(kubeconfigBase64)
	if err != nil {
		return fmt.Errorf("failed to decode kubeconfig: %v\n", err)

	}

	tmpDir, err := os.MkdirTemp("", "kubeconfig-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v\n", err)
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("failed to remove temp directory: %v\n", err)
		}
	}()

	kubeconfigPath := filepath.Join(tmpDir, "kubeconfig")
	err = os.WriteFile(kubeconfigPath, kubeconfigBytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to write kubeconfig: %v\n", err)
	}

	return fn(kubeconfigPath)
}
