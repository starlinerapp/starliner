package handler

import (
	"bufio"
	"io"
	"log"

	"google.golang.org/grpc"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
	"starliner.app/internal/provisioner/application"
)

type LogsHandler struct {
	v1.UnimplementedProvisioningLogServiceServer
	logsApplication *application.ClusterLogApplication
}

func NewLogsHandler(logsApplication *application.ClusterLogApplication) *LogsHandler {
	return &LogsHandler{
		logsApplication: logsApplication,
	}
}

func (lh *LogsHandler) StreamLogs(req *v1.StreamProvisioningLogRequest, stream grpc.ServerStreamingServer[v1.StreamProvisioningLogResponse]) error {
	rc, err := lh.logsApplication.StreamProvisioningLogs(stream.Context(), req.ProvisioningId)
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
		if err := stream.Send(&v1.StreamProvisioningLogResponse{Chunk: line}); err != nil {
			return err
		}
	}

	return scanner.Err()
}
