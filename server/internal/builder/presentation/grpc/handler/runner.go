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
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

const runnerTokenMetadataKey = "authorization"

type RunnerHandler struct {
	v1.UnimplementedRunnerSchedulerServiceServer
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

func (h *RunnerHandler) Connect(
	stream grpc.BidiStreamingServer[v1.RunnerMessage, v1.SchedulerMessage],
) error {
	token := tokenFromContext(stream.Context())

	runnerId, err := h.runnerApplication.ResolveRunnerId(stream.Context(), token)
	if err != nil {
		return err
	}

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return status.Errorf(codes.Unavailable, "receive runner message %v", err)
		}

		if msg == nil {
			continue
		}

		switch payload := msg.GetPayload().(type) {
		case *v1.RunnerMessage_Heartbeat:
			heartbeatInterval, err := h.heartbeatApplication.AcknowledgeHeartbeat(stream.Context(), runnerId)
			if err != nil {
				return status.Errorf(codes.Internal, "acknowledge heartbeat: %v", err)
			}

			err = stream.Send(&v1.SchedulerMessage{
				Payload: &v1.SchedulerMessage_HeartbeatAck{
					HeartbeatAck: &v1.HeartbeatAck{
						Sequence: payload.Heartbeat.GetSequence(),
						LeaseTtl: durationpb.New(heartbeatInterval),
					},
				},
			})
			if err != nil {
				return status.Errorf(codes.Unavailable, "send heartbeat ack: %v", err)
			}
		default:
			return status.Errorf(codes.InvalidArgument, "unsupported message type %T", payload)
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
