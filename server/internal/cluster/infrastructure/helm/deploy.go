package helm

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"helm.sh/helm/v4/pkg/action"
	v2 "helm.sh/helm/v4/pkg/chart/v2"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"starliner.app/internal/cluster/conf"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/infrastructure/shared/kubeconfig"
)

const (
	defaultStorage      = "secret"
	defaultWaitStrategy = "watcher"
)

type Deploy struct {
	config *conf.Config
}

func NewDeploy(config *conf.Config) port.Deploy {
	return &Deploy{
		config: config,
	}
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

	auth := base64.StdEncoding.EncodeToString(
		[]byte(args.ImageRegistryUsername + ":" + args.ImageRegistryPassword),
	)

	dockerConfig := map[string]interface{}{
		"auths": map[string]interface{}{
			args.ImageRegistryUrl: map[string]interface{}{
				"username": args.ImageRegistryUsername,
				"password": args.ImageRegistryPassword,
				"auth":     auth,
			},
		},
	}

	dockerConfigBytes, err := json.Marshal(dockerConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal docker config json: %w", err)
	}

	encodedDockerConfig := base64.StdEncoding.EncodeToString(dockerConfigBytes)

	values := map[string]interface{}{
		"image": map[string]interface{}{
			"repository": args.ImageName,
			"tag":        args.ImageTag,
		},
		"imagePullSecret": map[string]interface{}{
			"dockerconfigjson": encodedDockerConfig,
		},
		"port":       args.Port,
		"targetPort": args.Port,
		"env":        envValues,
	}

	if args.VolumeSizeMiB != nil && *args.VolumeSizeMiB > 0 {
		if args.VolumeMountPath == nil || *args.VolumeMountPath == "" {
			return fmt.Errorf("volume mount path is required when volume size is specified")
		}
		values["volume"] = map[string]interface{}{
			"sizeMiB":   *args.VolumeSizeMiB,
			"mountPath": *args.VolumeMountPath,
		}
		return installChart(StatefulSetChart, args.Namespace, args.ReleaseName, args.KubeconfigBase64, values)
	}

	return installChart(DeploymentChart, args.Namespace, args.ReleaseName, args.KubeconfigBase64, values)
}

func (d *Deploy) DeployPostgres(namespace string, releaseName string, kubeconfigBase64 string) error {
	values := map[string]interface{}{}
	return installChart(PostgresChart, namespace, releaseName, kubeconfigBase64, values)
}

func (d *Deploy) DeployExternalDNS(namespace string, releaseName string, kubeconfigBase64 string) error {
	values := map[string]interface{}{
		"cloudflare": map[string]interface{}{
			"apiKey": d.config.CFApiToken,
		},
		"external-dns": map[string]interface{}{
			"provider": map[string]interface{}{
				"name": "cloudflare",
			},
			"policy": "upsert-only",
			"env": []interface{}{
				map[string]interface{}{
					"name": "CF_API_TOKEN",
					"valueFrom": map[string]interface{}{
						"secretKeyRef": map[string]interface{}{
							"name": fmt.Sprintf("%s-cloudflare-api-key", releaseName),
							"key":  "apiKey",
						},
					},
				},
			},
		},
	}
	return installChart(ExternalDNSChart, namespace, releaseName, kubeconfigBase64, values)
}

func (d *Deploy) DeleteDeployment(namespace string, releaseName string, kubeconfigBase64 string) error {
	return uninstallRelease(namespace, releaseName, kubeconfigBase64)
}

func (d *Deploy) DeployIngress(args *port.DeployIngressArgs) error {
	values := buildIngressValues(args.Hosts, args.TLSEnabled)
	return installChart(IngressChart, args.Namespace, args.ReleaseName, args.KubeconfigBase64, values)
}

func buildIngressValues(hosts []port.IngressHost, tlsEnabled bool) map[string]interface{} {
	rules := make([]interface{}, 0, len(hosts))
	hostnames := make([]string, 0, len(hosts))

	for _, host := range hosts {
		hostnames = append(hostnames, host.Host)

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

	return map[string]interface{}{
		"ingress": map[string]interface{}{
			"rules":              rules,
			"hostnameAnnotation": strings.Join(hostnames, ","),
			"tls": map[string]interface{}{
				"enabled": tlsEnabled,
			},
		},
	}
}

func installChart(
	chartFS fs.FS,
	namespace string,
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

			cfg, err := newActionConfig(kubeconfigPath, namespace)
			if err != nil {
				return err
			}

			if values == nil {
				values = map[string]interface{}{}
			}

			if releaseExists(cfg, releaseName) {
				return upgradeChart(cfg, namespace, releaseName, ch, values)
			}
			return installNewChart(cfg, namespace, releaseName, ch, values)
		})
	})
}

func releaseExists(cfg *action.Configuration, releaseName string) bool {
	hist := action.NewHistory(cfg)
	hist.Max = 1
	releases, err := hist.Run(releaseName)
	return err == nil && len(releases) > 0
}

func installNewChart(cfg *action.Configuration, namespace string, releaseName string, ch *v2.Chart, values map[string]interface{}) error {
	install := action.NewInstall(cfg)
	install.ReleaseName = releaseName
	install.Namespace = namespace
	install.Timeout = 5 * time.Minute
	install.WaitStrategy = defaultWaitStrategy
	install.WaitForJobs = true
	install.CreateNamespace = true

	if _, err := install.Run(ch, values); err != nil {
		return fmt.Errorf("helm install %q failed: %w", releaseName, err)
	}
	return nil
}

func upgradeChart(cfg *action.Configuration, namespace string, releaseName string, ch *v2.Chart, values map[string]interface{}) error {
	upgrade := action.NewUpgrade(cfg)
	upgrade.Namespace = namespace
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

func uninstallRelease(namespace string, releaseName string, kubeconfigBase64 string) error {
	return kubeconfig.WithTempKubeConfig(kubeconfigBase64, func(kubeconfigPath string) error {
		cfg, err := newActionConfig(kubeconfigPath, namespace)
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
	flags.Namespace = &namespace

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
