package handler

import (
	"bufio"
	"io"
	"log"

	"google.golang.org/grpc"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/domain/value"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type LogsHandler struct {
	v1.UnimplementedLogsServiceServer
	logsApplication *application.LogsApplication
}

func NewLogsHandler(logsApplication *application.LogsApplication) *LogsHandler {
	return &LogsHandler{
		logsApplication: logsApplication,
	}
}

func (lh *LogsHandler) StreamLogs(req *v1.StreamLogsRequest, stream grpc.ServerStreamingServer[v1.StreamLogsResponse]) error {
	rc, err := lh.logsApplication.StreamDeploymentLogs(
		stream.Context(),
		logSourceFromProto(req.GetSource()),
		req.GetNamespace(),
		req.GetReleaseName(),
		req.GetKubeconfigBase64(),
	)
	if err != nil {
		return err
	}
	defer func(rc io.ReadCloser) {
		err := rc.Close()
		if err != nil {
			log.Printf("failed to close reader: %v", err)
		}
	}(rc)

	scanner := bufio.NewScanner(rc)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if err := stream.Send(&v1.StreamLogsResponse{Chunk: line}); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func logSourceFromProto(source v1.LogSource) value.LogSource {
	switch source {
	case v1.LogSource_LOG_SOURCE_INGRESS:
		return value.LogSourceIngress
	default:
		return value.LogSourceWorkload
	}
}
