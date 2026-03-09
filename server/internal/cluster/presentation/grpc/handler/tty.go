package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/domain/port"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type TtyHandler struct {
	v1.UnimplementedTTYServiceServer
	ttyApplication *application.TTYApplication
}

func NewTtyHandler(
	ttyApplication *application.TTYApplication,
) *TtyHandler {
	return &TtyHandler{
		ttyApplication: ttyApplication,
	}
}

// TODO: place orchestration logic in application layer
func (th *TtyHandler) OpenTTY(stream v1.TTYService_OpenTTYServer) error {
	ctx := stream.Context()

	firstReq, err := stream.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("missing initial tty session message")
		}
		return err
	}

	session, ok := firstReq.Payload.(*v1.OpenTTYRequest_Session)
	if !ok {
		return fmt.Errorf("first OpenTTY message must contain session payload")
	}

	sizeCh := make(chan port.TerminalSize, 1)
	defer close(sizeCh)

	stdin, stdout, err := th.ttyApplication.OpenTTY(
		ctx,
		session.Session.Namespace,
		session.Session.ReleaseName,
		session.Session.KubeconfigBase64,
		sizeCh,
	)
	if err != nil {
		return err
	}
	defer func(stdin io.WriteCloser) {
		err := stdin.Close()
		if err != nil {
			log.Printf("failed to close stdin: %v", err)
		}
	}(stdin)
	defer func(stdout io.ReadCloser) {
		err := stdout.Close()
		if err != nil {
			log.Printf("failed to close stdout: %v", err)
		}
	}(stdout)

	// Stream stdout back to the client
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				_ = stream.Send(&v1.OpenTTYResponse{Stdout: buf[:n]})
			}
			if err != nil {
				if err != io.EOF {
					_ = stream.Send(&v1.OpenTTYResponse{Stdout: []byte(fmt.Sprintf("error: %v", err))})
				}
				return
			}
		}
	}()

	// Receive stdin and resize events from the client
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream recv error: %w", err)
		}

		switch payload := msg.Payload.(type) {
		case *v1.OpenTTYRequest_Stdin:
			if _, err := stdin.Write(payload.Stdin); err != nil {
				return fmt.Errorf("stdin write error: %w", err)
			}
		case *v1.OpenTTYRequest_Size:
			select {
			case sizeCh <- port.TerminalSize{
				Columns: uint16(payload.Size.Cols),
				Rows:    uint16(payload.Size.Rows),
			}:
			default:
			}
		}
	}
}
