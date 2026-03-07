package port

import "context"

type GrpcClient interface {
	StreamLogs(
		ctx context.Context,
		namespace string,
		releaseName string,
		kubeconfigBase64 string,
	) error
}
