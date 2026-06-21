package handler

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"starliner.app/internal/builder/application"
	"starliner.app/internal/builder/domain/value"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

const runnerTokenMetadataKey = "authorization"

type RunnerHandler struct {
	v1.UnimplementedRunnerHeartbeatServiceServer
	runnerApplication    *application.RunnerApplication
	heartbeatApplication *application.HeartbeatApplication
}

func NewRunnerHandler(
	runnerApplication *application.RunnerApplication,
	heartbeatApplication *application.HeartbeatApplication,
) *RunnerHandler {
	return &RunnerHandler{
		runnerApplication:    runnerApplication,
		heartbeatApplication: heartbeatApplication,
	}
}

func (h *RunnerHandler) StreamHeartbeats(
	stream grpc.BidiStreamingServer[v1.HeartbeatMessage, v1.HeartbeatAckMessage],
) error {
	token := tokenFromContext(stream.Context())

	resolved, err := h.runnerApplication.ResolveRunner(stream.Context(), token)
	if err != nil {
		return err
	}

	runnerId := resolved.Id
	organizationId := resolved.OrganizationId

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return status.Errorf(codes.Unavailable, "receive heartbeat: %v", err)
		}

		if msg == nil {
			continue
		}

		heartbeat := msg.GetHeartbeat()
		if heartbeat == nil {
			return status.Error(codes.InvalidArgument, "missing heartbeat payload")
		}

		heartbeatInterval, err := h.heartbeatApplication.AcknowledgeHeartbeat(
			stream.Context(),
			runnerId,
			organizationId,
			heartbeat.GetMaxConcurrentJobs(),
			heartbeat.GetActiveJobs(),
		)
		if err != nil {
			if errors.Is(err, value.ErrRunnerDeleted) {
				return status.Error(codes.FailedPrecondition, "runner deleted")
			}
			if errors.Is(err, value.ErrInvalidRunnerCapacity) {
				return status.Error(codes.InvalidArgument, err.Error())
			}
			return status.Errorf(codes.Internal, "acknowledge heartbeat: %v", err)
		}

		if err := stream.Send(&v1.HeartbeatAckMessage{
			Payload: &v1.HeartbeatAckMessage_HeartbeatAck{
				HeartbeatAck: &v1.HeartbeatAck{
					Sequence: heartbeat.GetSequence(),
					LeaseTtl: durationpb.New(heartbeatInterval),
				},
			},
		}); err != nil {
			return status.Errorf(codes.Unavailable, "send heartbeat ack: %v", err)
		}
	}
}

func tokenFromContext(ctx context.Context) string {
	values := metadata.ValueFromIncomingContext(ctx, runnerTokenMetadataKey)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}
