package handler

import (
	"bufio"
	"io"
	"log"

	"google.golang.org/grpc"
	"starliner.app/internal/builder/application"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type BuildLogHandler struct {
	v1.UnimplementedBuildLogServiceServer
	buildLogApplication *application.BuildLogApplication
}

func NewBuildLogHandler(app *application.BuildLogApplication) *BuildLogHandler {
	return &BuildLogHandler{buildLogApplication: app}
}

func (b *BuildLogHandler) StreamBuildLogs(
	req *v1.StreamBuildLogRequest,
	stream grpc.ServerStreamingServer[v1.StreamBuildLogResponse],
) error {
	rc, err := b.buildLogApplication.StreamBuildLogs(stream.Context(), req.GetBuildId())
	if err != nil {
		return err
	}
	defer func(rc io.ReadCloser) {
		if err := rc.Close(); err != nil {
			log.Printf("failed to close build log reader: %v", err)
		}
	}(rc)

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Bytes()
		if err := stream.Send(&v1.StreamBuildLogResponse{Chunk: line}); err != nil {
			return err
		}
	}

	return scanner.Err()
}
