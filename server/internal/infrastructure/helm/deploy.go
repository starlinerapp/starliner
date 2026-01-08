package helm

import (
	"fmt"
	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	"io/fs"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"path/filepath"
	"starliner.app/internal/domain/port"
)

type Deploy struct {
}

func NewDeploy() port.Deploy {
	return &Deploy{}
}

func (d *Deploy) DeployNginx(ip string, kubeconfigPath string) error {
	tmpDir, err := os.MkdirTemp("", "helm-chart-*")
	if err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("failed to remove temp directory: %v\n", err)
		}
	}()

	err = fs.WalkDir(NginxChart, "nginx", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("nginx", path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(tmpDir, relPath)
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := NginxChart.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, 0644)
	})
	if err != nil {
		return err
	}

	chart, err := loader.Load(tmpDir)
	if err != nil {
		return err
	}

	configFlags := genericclioptions.NewConfigFlags(false)
	configFlags.KubeConfig = &kubeconfigPath

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(
		configFlags,
		"default",
		"secret",
	); err != nil {
		return err
	}

	install := action.NewInstall(actionConfig)
	install.ReleaseName = "test-release"
	install.Namespace = "default"
	install.WaitStrategy = "watcher"

	fmt.Printf("Installing helm chart")
	_, err = install.Run(chart, map[string]interface{}{
		"ingress": map[string]interface{}{
			"host": fmt.Sprintf("%s.nip.io", ip),
		},
	})

	return nil
}
