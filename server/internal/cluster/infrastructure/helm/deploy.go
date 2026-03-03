package helm

import (
	"fmt"
	"helm.sh/helm/v4/pkg/action"
	v2 "helm.sh/helm/v4/pkg/chart/v2"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	"io/fs"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"path/filepath"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/infrastructure/shared/kubeconfig"
	"time"
)

const (
	defaultNamespace    = "default"
	defaultStorage      = "secret"
	defaultWaitStrategy = "watcher"
)

type Deploy struct {
}

func NewDeploy() port.Deploy {
	return &Deploy{}
}

func (d *Deploy) DeployImage(
	args *port.DeployImageArgs,
) error {
	envValues := make([]map[string]interface{}, 0, len(args.EnvVars))

	for _, e := range args.EnvVars {
		envValues = append(envValues, map[string]interface{}{
			"name":  e.Name,
			"value": e.Value,
		})
	}

	return installChart(ApplicationChart, args.ReleaseName, args.KubeconfigBase64, map[string]interface{}{
		"image": map[string]interface{}{
			"repository": args.ImageRepository,
			"tag":        args.ImageTag,
		},
		"port":       args.Port,
		"targetPort": args.Port,
		"env":        envValues,
	})
}

func (d *Deploy) DeployCloudNativePg(releaseName string, kubeconfigBase64 string) error {
	return installChart(CloudNativePgChart, releaseName, kubeconfigBase64, nil)
}

func (d *Deploy) DeployPostgres(releaseName string, kubeconfigBase64 string) error {
	return installChart(PostgresChart, releaseName, kubeconfigBase64, nil)
}

func (d *Deploy) DeleteDeployment(releaseName string, kubeconfigBase64 string) error {
	return uninstallRelease(releaseName, kubeconfigBase64)
}

func (d *Deploy) DeployIngress(args *port.DeployIngressArgs) error {
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
	return installChart(IngressChart, args.ReleaseName, args.KubeconfigBase64, values)
}

func installChart(
	chartFS fs.FS,
	releaseName string,
	kubeconfigBase64 string,
	values map[string]interface{},
) error {
	return kubeconfig.WithTempKubeConfig(kubeconfigBase64, func(kubeconfigPath string) error {
		return withTempChartDir(chartFS, func(chartDir string) error {
			ch, err := loader.Load(chartDir)
			if err != nil {
				return fmt.Errorf("failed to load chart from %q: %w", chartDir, err)
			}

			cfg, err := newActionConfig(kubeconfigPath, defaultNamespace)
			if err != nil {
				return err
			}

			if values == nil {
				values = map[string]interface{}{}
			}

			if releaseExists(cfg, releaseName) {
				return upgradeChart(cfg, releaseName, ch, values)
			}
			return installNewChart(cfg, releaseName, ch, values)
		})
	})
}

func releaseExists(cfg *action.Configuration, releaseName string) bool {
	hist := action.NewHistory(cfg)
	hist.Max = 1
	releases, err := hist.Run(releaseName)
	return err == nil && len(releases) > 0
}

func installNewChart(cfg *action.Configuration, releaseName string, ch *v2.Chart, values map[string]interface{}) error {
	install := action.NewInstall(cfg)
	install.ReleaseName = releaseName
	install.Namespace = defaultNamespace
	install.Timeout = 5 * time.Minute
	install.WaitStrategy = defaultWaitStrategy
	install.WaitForJobs = true
	install.CreateNamespace = true

	if _, err := install.Run(ch, values); err != nil {
		return fmt.Errorf("helm install %q failed: %w", releaseName, err)
	}
	return nil
}

func upgradeChart(cfg *action.Configuration, releaseName string, ch *v2.Chart, values map[string]interface{}) error {
	upgrade := action.NewUpgrade(cfg)
	upgrade.Namespace = defaultNamespace
	upgrade.Timeout = 5 * time.Minute
	upgrade.WaitStrategy = defaultWaitStrategy
	upgrade.WaitForJobs = true
	upgrade.RollbackOnFailure = true
	upgrade.CleanupOnFail = true

	if _, err := upgrade.Run(releaseName, ch, values); err != nil {
		return fmt.Errorf("helm upgrade %q failed: %w", releaseName, err)
	}
	return nil
}

func uninstallRelease(releaseName string, kubeconfigBase64 string) error {
	return kubeconfig.WithTempKubeConfig(kubeconfigBase64, func(kubeconfigPath string) error {
		cfg, err := newActionConfig(kubeconfigPath, defaultNamespace)
		if err != nil {
			return err
		}

		uninstall := action.NewUninstall(cfg)
		uninstall.WaitStrategy = defaultWaitStrategy

		if _, err = uninstall.Run(releaseName); err != nil {
			return fmt.Errorf("helm uninstall %q failed: %w", releaseName, err)
		}
		return nil
	})
}

func newActionConfig(kubeconfigPath string, namespace string) (*action.Configuration, error) {
	flags := genericclioptions.NewConfigFlags(false)
	flags.KubeConfig = &kubeconfigPath

	cfg := new(action.Configuration)
	if err := cfg.Init(flags, namespace, defaultStorage); err != nil {
		return nil, fmt.Errorf("failed to init helm action config: %w", err)
	}
	return cfg, nil
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
