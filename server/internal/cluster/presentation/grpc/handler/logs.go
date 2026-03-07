package handler

import (
	"bufio"
	"google.golang.org/grpc"
	"io"
	"log"
	"starliner.app/internal/cluster/application"
	v2 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type LogsHandler struct {
	v2.UnimplementedLogsServiceServer
	logsApplication *application.LogsApplication
}

func NewLogsHandler(logsApplication *application.LogsApplication) *LogsHandler {
	return &LogsHandler{
		logsApplication: logsApplication,
	}
}

func (lh *LogsHandler) StreamLogs(req *v2.StreamLogsRequest, stream grpc.ServerStreamingServer[v2.StreamLogsResponse]) error {
	rc, err := lh.logsApplication.StreamDeploymentLogs(stream.Context(), req.Namespace, req.ReleaseName, req.KubeconfigBase64)
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
	for scanner.Scan() {
		line := scanner.Bytes()
		if err := stream.Send(&v2.StreamLogsResponse{Chunk: line}); err != nil {
			return err
		}
	}

	return scanner.Err()
}
