package helm

import (
	"encoding/base64"
	"fmt"
	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	"io/fs"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/cluster/domain/port"
)

type Deploy struct {
}

func NewDeploy() port.Deploy {
	return &Deploy{}
}

func (d *Deploy) DeployCloudNativePg(releaseName string, kubeconfigBase64 string) error {
	return withTempKubeConfig(kubeconfigBase64, func(kubeconfigPath string) error {
		return withTempChartDir(CloudNativePgChart, func(tmpDir string) error {
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

			install.ReleaseName = releaseName
			install.Namespace = "default"
			install.WaitStrategy = "watcher"

			log.Println("Installing helm chart " + releaseName + "...")
			_, err = install.Run(chart, map[string]interface{}{})
			if err != nil {
				return err
			}

			return nil
		})
	})
}

func (d *Deploy) DeployPostgres(releaseName string, kubeconfigBase64 string) error {
	return withTempKubeConfig(kubeconfigBase64, func(kubeconfigPath string) error {
		return withTempChartDir(PostgresChart, func(tmpDir string) error {
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

			install.ReleaseName = releaseName
			install.Namespace = "default"
			install.WaitStrategy = "watcher"

			log.Println("Installing helm chart " + releaseName + "...")
			_, err = install.Run(chart, map[string]interface{}{})
			if err != nil {
				return err
			}

			return nil
		})
	})
}

func (d *Deploy) DeletePostgres(releaseName string, kubeconfigBase64 string) error {
	return withTempKubeConfig(kubeconfigBase64, func(kubeconfigPath string) error {
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

		uninstall := action.NewUninstall(actionConfig)
		uninstall.WaitStrategy = "watcher"

		_, err := uninstall.Run(releaseName)
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *Deploy) DeployIngress(args *port.DeployIngressArgs) error {
	return withTempKubeConfig(args.KubeconfigBase64, func(kubeconfigPath string) error {
		return withTempChartDir(IngressChart, func(tmpDir string) error {
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

			install.ReleaseName = args.ReleaseName
			install.Namespace = "default"
			install.WaitStrategy = "watcher"

			log.Println("Installing helm chart " + args.ReleaseName + "...")

			rules := make([]interface{}, 0, len(args.Hosts))
			for _, host := range args.Hosts {
				paths := make([]interface{}, 0, len(host.Paths))
				for _, p := range host.Paths {
					paths = append(paths, map[string]interface{}{
						"path":     p.Path,
						"pathType": p.PathType,
						"service": map[string]interface{}{
							"name": p.ServiceName,
							"port": map[string]interface{}{
								"number": p.ServicePort,
							},
						},
					})
				}

				rules = append(rules, map[string]interface{}{
					"host":  host.Host,
					"paths": paths,
				})
			}

			values := map[string]interface{}{
				"ingress": map[string]interface{}{
					"rules": rules,
				},
			}
			_, err = install.Run(chart, values)
			if err != nil {
				return err
			}

			return nil
		})
	})
}

func withTempKubeConfig(kubeconfigBase64 string, fn func(path string) error) error {
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

func withTempChartDir(chartFS fs.FS, fn func(path string) error) error {
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

	entries, err := fs.ReadDir(chartFS, ".")
	if err != nil {
		return err
	}
	if len(entries) != 1 || !entries[0].IsDir() {
		return fmt.Errorf("expected one top-level directory in embedded FS")
	}

	root := entries[0].Name()
	err = fs.WalkDir(chartFS, root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(tmpDir, relPath)
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := fs.ReadFile(chartFS, path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0644)
	})
	if err != nil {
		return err
	}

	return fn(tmpDir)
}
