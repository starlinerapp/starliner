package application

import (
	"context"
	"io"

	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/domain/value"
)

type LogsApplication struct {
	logs port.Logs
}

func NewLogsApplication(logs port.Logs) *LogsApplication {
	return &LogsApplication{logs: logs}
}

func (la *LogsApplication) StreamDeploymentLogs(
	ctx context.Context,
	source value.LogSource,
	environmentNamespace string,
	releaseName string,
	kubeconfigBase64 string,
) (io.ReadCloser, error) {
	return la.logs.StreamLogs(ctx, source, environmentNamespace, releaseName, kubeconfigBase64)
}
