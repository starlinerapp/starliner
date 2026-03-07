package application

import (
	"context"
	"io"
	"starliner.app/internal/cluster/domain/port"
)

type LogsApplication struct {
	logs port.Logs
}

func NewLogsApplication(logs port.Logs) *LogsApplication {
	return &LogsApplication{logs: logs}
}

func (la *LogsApplication) StreamDeploymentLogs(ctx context.Context, namespace string, releaseName string, kubeconfigBase64 string) (io.ReadCloser, error) {
	return la.logs.StreamLogs(ctx, namespace, releaseName, kubeconfigBase64)
}
